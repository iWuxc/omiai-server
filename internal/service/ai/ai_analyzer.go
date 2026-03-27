package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"omiai-server/internal/conf"
	"strings"
	"time"

	"github.com/iWuxc/go-wit/log"
)

type AIAnalyzer struct {
	provider LLMProvider
}

type LLMProvider interface {
	Call(prompt string) (string, error)
	Name() string
}

func NewAIAnalyzer() *AIAnalyzer {
	return &AIAnalyzer{
		provider: getAIProvider(),
	}
}

func getAIProvider() LLMProvider {
	cfg := conf.GetConfig()
	if cfg == nil {
		log.Warn("AI config is nil")
		return &MockAIProvider{}
	}
	if cfg.LLM == nil {
		log.Warn("AI LLM config is nil")
		return &MockAIProvider{}
	}
	log.Infof("LLM config: provider=%s, volcano_api_key=%s, model=%s",
		cfg.LLM.Provider, cfg.LLM.VolcanoEngine.APIKey, cfg.LLM.VolcanoEngine.Model)

	if cfg.LLM.VolcanoEngine != nil && cfg.LLM.VolcanoEngine.APIKey != "" {
		return &VolcanoAIProvider{
			APIKey:   cfg.LLM.VolcanoEngine.APIKey,
			Model:    cfg.LLM.VolcanoEngine.Model,
			Endpoint: cfg.LLM.VolcanoEngine.Endpoint,
		}
	}

	return &MockAIProvider{}
}

type VolcanoAIProvider struct {
	APIKey   string
	Model    string
	Endpoint string
}

func (p *VolcanoAIProvider) Name() string {
	return "volcano"
}

func (p *VolcanoAIProvider) Call(prompt string) (string, error) {
	endpoint := p.Endpoint
	if endpoint == "" {
		endpoint = "https://ark.cn-beijing.volces.com/api/coding/v3/chat/completions"
	}

	requestBody := map[string]interface{}{
		"model": p.Model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"temperature": 0.7,
		"max_tokens":  2000,
	}

	jsonData, _ := json.Marshal(requestBody)

	client := &http.Client{Timeout: 180 * time.Second}
	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.APIKey)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	log.Infof("Volcano AI raw response: %s", string(body))

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if result.Error.Message != "" {
		return "", fmt.Errorf("AI API error: %s", result.Error.Message)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	return result.Choices[0].Message.Content, nil
}

type MockAIProvider struct{}

func (p *MockAIProvider) Name() string {
	return "mock"
}

func (p *MockAIProvider) Call(prompt string) (string, error) {
	return `{"summary": "Mock summary", "extracted_tags": ["mock"], "emotion": "中立"}`, nil
}

type MatchAnalysisRequest struct {
	ClientA *ClientProfile `json:"client_a"`
	ClientB *ClientProfile `json:"client_b"`
}

type ClientProfile struct {
	Name                string `json:"name"`
	Gender              string `json:"gender"`
	Age                 int    `json:"age"`
	Height              int    `json:"height"`
	Education           string `json:"education"`
	Income              int    `json:"income"`
	Profession          string `json:"profession"`
	MaritalStatus       string `json:"marital_status"`
	HouseStatus         string `json:"house_status"`
	CarStatus           string `json:"car_status"`
	Address             string `json:"address"`
	FamilyDescription   string `json:"family_description"`
	PartnerRequirements string `json:"partner_requirements"`
	Remark              string `json:"remark"`
	Tags                string `json:"tags"`
}

type MatchAnalysisResult struct {
	OverallScore       int                `json:"overall_score"`
	Level              string             `json:"level"`
	HardConditions     *ConditionAnalysis `json:"hard_conditions"`
	SoftConditions     *ConditionAnalysis `json:"soft_conditions"`
	RiskPoints         []string           `json:"risk_points"`
	Advantages         []string           `json:"advantages"`
	Suggestions        string             `json:"suggestions"`
	IceBreakerTopics   []string           `json:"ice_breaker_topics"`
	SuccessProbability string             `json:"success_probability"`
}

type ConditionAnalysis struct {
	Score       int      `json:"score"`
	MatchLevel  string   `json:"match_level"`
	Analysis    string   `json:"analysis"`
	Suggestions []string `json:"suggestions"`
}

type ChatSummaryResult struct {
	Summary       string   `json:"summary"`
	ExtractedTags []string `json:"extracted_tags"`
	Emotion       string   `json:"emotion"`
}

func (a *AIAnalyzer) AnalyzeChatSummary(chatContent string) (*ChatSummaryResult, error) {
	prompt := buildChatSummaryPrompt(chatContent)
	response, err := a.provider.Call(prompt)
	if err != nil {
		return nil, err
	}

	var result ChatSummaryResult
	jsonStr := cleanJSONResponse(response)
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("解析AI响应失败: %v, raw: %s", err, response)
	}
	return &result, nil
}

func buildChatSummaryPrompt(chatContent string) string {
	return fmt.Sprintf(`【角色设定】
你是一位专业的婚恋分析师，擅长从客户的日常聊天记录中敏锐地捕捉他们的性格偏好和真实需求。

【任务】
请阅读以下脱敏聊天记录，提取客户的隐性偏好标签，并生成一段简短的跟进纪要。

【聊天记录】
%s

【输出要求】
必须严格按以下JSON格式输出，不要添加任何markdown代码块标记，不要有任何解释性文字：

{
  "summary": "客户今天表达了对上次相亲对象的不满，主要原因是对方迟到。客户比较看重时间观念。",
  "extracted_tags": ["看重时间观念", "讨厌迟到", "直性子"],
  "emotion": "消极"
}

【重要】只输出JSON，不要任何其他文字！`, chatContent)
}

func (a *AIAnalyzer) AnalyzeMatch(clientA, clientB *ClientProfile) (*MatchAnalysisResult, error) {
	prompt := buildMatchAnalysisPrompt(clientA, clientB)

	response, err := a.provider.Call(prompt)
	if err != nil {
		return nil, err
	}

	return parseAIResponse(response)
}

func buildMatchAnalysisPrompt(clientA, clientB *ClientProfile) string {
	return fmt.Sprintf(`【角色设定】
你是一位资深的婚恋顾问，拥有20年的红娘经验，擅长分析客户的匹配度。

【任务】
分析以下两位客户的匹配度，给出专业的分析和建议。

【客户A资料】
姓名：%s
性别：%s
年龄：%d岁
身高：%dcm
学历：%s
月收入：%d元
职业：%s
婚姻状况：%s
房产情况：%s
车辆情况：%s
居住地址：%s
家庭情况：%s
择偶要求：%s
红娘备注：%s
标签：%s

【客户B资料】
姓名：%s
性别：%s
年龄：%d岁
身高：%dcm
学历：%s
月收入：%d元
职业：%s
婚姻状况：%s
房产情况：%s
车辆情况：%s
居住地址：%s
家庭情况：%s
择偶要求：%s
红娘备注：%s
标签：%s

【输出要求】
必须严格按以下JSON格式输出，不要添加任何markdown代码块标记，不要有任何解释性文字：

{
  "overall_score": 85,
  "level": "良好匹配",
  "hard_conditions": {
    "score": 80,
    "match_level": "高度匹配",
    "analysis": "年龄、身高、学历等硬性条件匹配度分析",
    "suggestions": ["建议1", "建议2"]
  },
  "soft_conditions": {
    "score": 75,
    "match_level": "基本匹配",
    "analysis": "性格、兴趣、价值观等软性条件匹配度分析",
    "suggestions": ["建议1"]
  },
  "risk_points": ["风险点1", "风险点2"],
  "advantages": ["优势1", "优势2"],
  "suggestions": "总体建议和行动方案",
  "ice_breaker_topics": ["话题1", "话题2", "话题3"],
  "success_probability": "根据分析，预测这对客户的成功概率为70%%-80%%，建议优先推荐。"
}

【字段说明】
- overall_score: 总体匹配分数，0-100的整数
- level: 匹配等级，可选：完美匹配、良好匹配、一般匹配、不太匹配
- hard_conditions.score: 硬性条件得分，0-100的整数
- soft_conditions.score: 软性条件得分，0-100的整数
- risk_points: 风险点数组，没有则填[]
- advantages: 优势点数组，没有则填[]
- suggestions: 总体建议，不超过200字
- ice_breaker_topics: 破冰话题数组，3-5个话题
- success_probability: 成功概率预测描述

【重要】只输出JSON，不要任何其他文字！`,
		clientA.Name, clientA.Gender, clientA.Age, clientA.Height, clientA.Education,
		clientA.Income, clientA.Profession, clientA.MaritalStatus, clientA.HouseStatus,
		clientA.CarStatus, clientA.Address, clientA.FamilyDescription,
		clientA.PartnerRequirements, clientA.Remark, clientA.Tags,
		clientB.Name, clientB.Gender, clientB.Age, clientB.Height, clientB.Education,
		clientB.Income, clientB.Profession, clientB.MaritalStatus, clientB.HouseStatus,
		clientB.CarStatus, clientB.Address, clientB.FamilyDescription,
		clientB.PartnerRequirements, clientB.Remark, clientB.Tags,
	)
}

func parseAIResponse(response string) (*MatchAnalysisResult, error) {
	cleaned := cleanJSONResponse(response)
	log.Infof("AI cleaned response: %s", cleaned)

	var result MatchAnalysisResult
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		return nil, fmt.Errorf("parse AI response failed: %v", err)
	}

	return &result, nil
}

func cleanJSONResponse(response string) string {
	response = strings.TrimSpace(response)
	response = strings.ReplaceAll(response, "```json", "")
	response = strings.ReplaceAll(response, "```JSON", "")
	response = strings.ReplaceAll(response, "```", "")
	response = strings.TrimPrefix(response, "json")

	startIdx := strings.Index(response, "{")
	if startIdx == -1 {
		startIdx = strings.Index(response, "[")
	}
	if startIdx > 0 {
		response = response[startIdx:]
	}

	endIdx := strings.LastIndex(response, "}")
	if endIdx == -1 {
		endIdx = strings.LastIndex(response, "]")
	}
	if endIdx > 0 && endIdx < len(response)-1 {
		response = response[:endIdx+1]
	}

	return strings.TrimSpace(response)
}

func (a *AIAnalyzer) GenerateIceBreaker(clientA, clientB *ClientProfile) ([]string, error) {
	prompt := fmt.Sprintf(`作为资深红娘，为以下两位客户推荐3-5个破冰话题，帮助他们初次交流时打破尴尬。

客户A：%s，%d岁，职业：%s，兴趣标签：%s
客户B：%s，%d岁，职业：%s，兴趣标签：%s

请输出JSON数组格式：["话题1", "话题2", "话题3"]`,
		clientA.Name, clientA.Age, clientA.Profession, clientA.Tags,
		clientB.Name, clientB.Age, clientB.Profession, clientB.Tags,
	)

	response, err := a.provider.Call(prompt)
	if err != nil {
		return nil, err
	}

	var topics []string
	if err := json.Unmarshal([]byte(cleanJSONResponse(response)), &topics); err != nil {
		return []string{
			"最近工作怎么样？",
			"平时周末喜欢做什么？",
			"对另一半有什么期待？",
		}, nil
	}

	return topics, nil
}
