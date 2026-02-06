package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"omiai-server/internal/conf"
	"strings"
)

// AIAnalyzer AI分析服务
type AIAnalyzer struct {
	apiKey string
	model  string
}

// NewAIAnalyzer 创建AI分析器
func NewAIAnalyzer() *AIAnalyzer {
	cfg := conf.GetConfig().ZhipuAI
	return &AIAnalyzer{
		apiKey: cfg.APIKey,
		model:  cfg.Model,
	}
}

// MatchAnalysisRequest AI匹配分析请求
type MatchAnalysisRequest struct {
	ClientA *ClientProfile `json:"client_a"`
	ClientB *ClientProfile `json:"client_b"`
}

// ClientProfile 客户档案（用于AI分析）
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

// MatchAnalysisResult AI匹配分析结果
type MatchAnalysisResult struct {
	OverallScore       int                `json:"overall_score"`       // 总体匹配度 0-100
	Level              string             `json:"level"`               // 匹配等级
	HardConditions     *ConditionAnalysis `json:"hard_conditions"`     // 硬性条件分析
	SoftConditions     *ConditionAnalysis `json:"soft_conditions"`     // 软性条件分析
	RiskPoints         []string           `json:"risk_points"`         // 风险点
	Advantages         []string           `json:"advantages"`          // 优势点
	Suggestions        string             `json:"suggestions"`         // 建议
	IceBreakerTopics   []string           `json:"ice_breaker_topics"`  // 破冰话题
	SuccessProbability string             `json:"success_probability"` // 成功概率预测
}

// ConditionAnalysis 条件分析
type ConditionAnalysis struct {
	Score       int      `json:"score"`       // 单项得分 0-100
	MatchLevel  string   `json:"match_level"` // 匹配程度
	Analysis    string   `json:"analysis"`    // 分析说明
	Suggestions []string `json:"suggestions"` // 改进建议
}

// AnalyzeMatch 分析两个客户的匹配度
func (a *AIAnalyzer) AnalyzeMatch(clientA, clientB *ClientProfile) (*MatchAnalysisResult, error) {
	prompt := buildMatchAnalysisPrompt(clientA, clientB)

	response, err := a.callZhipuAI(prompt)
	if err != nil {
		return nil, err
	}

	return parseAIResponse(response)
}

// buildMatchAnalysisPrompt 构建匹配分析Prompt
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
  "success_probability": "根据分析，预测这对客户的成功概率为70%-80%，建议优先推荐。"
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

// callZhipuAI 调用智谱AI API
func (a *AIAnalyzer) callZhipuAI(prompt string) (string, error) {
	url := "https://open.bigmodel.cn/api/paas/v4/chat/completions"

	requestBody := map[string]interface{}{
		"model": a.model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"temperature": 0.7,
		"max_tokens":  2000,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

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

// parseAIResponse 解析AI返回的JSON
func parseAIResponse(response string) (*MatchAnalysisResult, error) {
	// 清理可能的Markdown代码块标记
	response = cleanJSONResponse(response)

	var result MatchAnalysisResult
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return nil, fmt.Errorf("parse AI response failed: %v", err)
	}

	return &result, nil
}

// cleanJSONResponse 清理AI返回的JSON字符串
func cleanJSONResponse(response string) string {
	response = strings.TrimSpace(response)

	// 移除Markdown代码块标记（多种格式）
	response = strings.ReplaceAll(response, "```json", "")
	response = strings.ReplaceAll(response, "```JSON", "")
	response = strings.ReplaceAll(response, "```", "")

	// 移除可能的 "json" 前缀
	response = strings.TrimPrefix(response, "json")

	// 找到JSON开始的位置（第一个 { 或 [）
	startIdx := strings.Index(response, "{")
	if startIdx == -1 {
		startIdx = strings.Index(response, "[")
	}
	if startIdx > 0 {
		response = response[startIdx:]
	}

	// 找到JSON结束的位置（最后一个 } 或 ]）
	endIdx := strings.LastIndex(response, "}")
	if endIdx == -1 {
		endIdx = strings.LastIndex(response, "]")
	}
	if endIdx > 0 && endIdx < len(response)-1 {
		response = response[:endIdx+1]
	}

	return strings.TrimSpace(response)
}

// GenerateIceBreaker 生成破冰话题
func (a *AIAnalyzer) GenerateIceBreaker(clientA, clientB *ClientProfile) ([]string, error) {
	prompt := fmt.Sprintf(`作为资深红娘，为以下两位客户推荐3-5个破冰话题，帮助他们初次交流时打破尴尬。

客户A：%s，%d岁，职业：%s，兴趣标签：%s
客户B：%s，%d岁，职业：%s，兴趣标签：%s

请输出JSON数组格式：["话题1", "话题2", "话题3"]`,
		clientA.Name, clientA.Age, clientA.Profession, clientA.Tags,
		clientB.Name, clientB.Age, clientB.Profession, clientB.Tags,
	)

	response, err := a.callZhipuAI(prompt)
	if err != nil {
		return nil, err
	}

	var topics []string
	if err := json.Unmarshal([]byte(cleanJSONResponse(response)), &topics); err != nil {
		// 如果解析失败，返回默认话题
		return []string{
			"最近工作怎么样？",
			"平时周末喜欢做什么？",
			"对另一半有什么期待？",
		}, nil
	}

	return topics, nil
}
