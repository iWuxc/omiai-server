package chat_parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"omiai-server/internal/conf"
	"strings"
	"time"

	"github.com/iWuxc/go-wit/log"
)

// ImportRecord 临时导入记录结构
type ImportRecord struct {
	Name                string `json:"name"`
	Gender              int8   `json:"gender"` // 1男 2女
	Phone               string `json:"phone"`
	Birthday            string `json:"birthday"`
	Height              int    `json:"height"`
	Weight              int    `json:"weight"`
	Education           int8   `json:"education"`
	MaritalStatus       int8   `json:"marital_status"`
	Income              int    `json:"income"`
	Address             string `json:"address"`
	Profession          string `json:"profession"`
	HouseStatus         int8   `json:"house_status"`
	CarStatus           int8   `json:"car_status"`
	PartnerRequirements string `json:"partner_requirements"` // 原文保留
	RawText             string `json:"raw_text"`             // 原始文本片段
	ParseStatus         string `json:"parse_status"`         // success, warning, error
	ErrorMsg            string `json:"error_msg"`
}

// ChatParser 解析器
type ChatParser struct {
	Records []ImportRecord
	APIKey  string
	Model   string
}

func NewChatParser() *ChatParser {
	// 默认配置
	apiKey := ""
	model := "glm-4-flash"

	// 从配置读取
	if c := conf.GetConfig(); c != nil && c.ZhipuAI != nil {
		if c.ZhipuAI.APIKey != "" {
			apiKey = c.ZhipuAI.APIKey
		}
		if c.ZhipuAI.Model != "" {
			model = c.ZhipuAI.Model
		}
	}

	fmt.Println("大模型key", apiKey, model)

	return &ChatParser{
		Records: make([]ImportRecord, 0),
		APIKey:  apiKey,
		Model:   model,
	}
}

// GLM Request/Response Structures
type GLMRequest struct {
	Model    string       `json:"model"`
	Messages []GLMMessage `json:"messages"`
}

type GLMMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GLMResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// Parse 解析文本块 (Using GLM)
func (p *ChatParser) Parse(content string) ([]ImportRecord, error) {
	if p.APIKey == "" {
		// Fallback to Mock if no key
		return p.parseMock(content), nil
	}

	prompt := `
你是一个专业的数据提取助手。请从以下聊天记录中提取客户资料，并严格输出为 JSON 数组格式。
忽略无关的闲聊，只提取包含客户信息的条目。

字段定义：
- name (string): 姓名
- gender (int): 1=男, 2=女
- phone (string): 手机号 (11位数字)
- birthday (string): 格式 YYYY-MM
- height (int): cm
- weight (int): kg
- education (int): 1=高中, 2=大专, 3=本科, 4=硕士, 5=博士
- marital_status (int): 1=未婚, 2=已婚, 3=离异, 4=丧偶
- income (int): 月收入(元)
- address (string): 居住地
- profession (string): 职业
- house_status (int): 1=无房, 2=有房
- car_status (int): 1=无车, 2=有车
- partner_requirements (string): 择偶要求摘要

请直接返回 JSON 数组，不要包含 Markdown 标记。

待处理文本：
"""
%s
"""
`
	fullPrompt := fmt.Sprintf(prompt, content)

	reqBody := GLMRequest{
		Model: p.Model,
		Messages: []GLMMessage{
			{Role: "user", Content: fullPrompt},
		},
	}

	jsonBody, _ := json.Marshal(reqBody)

	// Increase timeout to 180s for slow model responses
	client := &http.Client{Timeout: 180 * time.Second}
	req, _ := http.NewRequest("POST", "https://open.bigmodel.cn/api/paas/v4/chat/completions", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.APIKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("GLM API request failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var glmResp GLMResponse
	if err := json.NewDecoder(resp.Body).Decode(&glmResp); err != nil {
		log.Errorf("GLM API response decode failed: %v", err)
		return nil, err
	}

	if len(glmResp.Choices) == 0 {
		return nil, fmt.Errorf("empty response from GLM")
	}

	contentStr := glmResp.Choices[0].Message.Content
	// Clean markdown code blocks if present
	contentStr = strings.TrimPrefix(contentStr, "```json")
	contentStr = strings.TrimPrefix(contentStr, "```")
	contentStr = strings.TrimSuffix(contentStr, "```")
	contentStr = strings.TrimSpace(contentStr)

	var records []ImportRecord
	if err := json.Unmarshal([]byte(contentStr), &records); err != nil {
		log.Errorf("JSON parse failed: %v, content: %s", err, contentStr)
		return nil, fmt.Errorf("failed to parse JSON from GLM")
	}

	// Post-processing
	for i := range records {
		records[i].ParseStatus = "success"
		if records[i].Name == "" {
			records[i].Name = "未知用户"
			records[i].ParseStatus = "warning"
		}
	}

	return records, nil
}

// Mock implementation for testing without API Key
func (p *ChatParser) parseMock(content string) []ImportRecord {
	return []ImportRecord{
		{
			Name:        "张三(Mock)",
			Gender:      1,
			Phone:       "13800138000",
			Birthday:    "1990-01",
			Height:      175,
			Education:   3,
			ParseStatus: "success",
			RawText:     "Mock Data",
		},
	}
}
