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

type ImportRecord struct {
	Name                string `json:"name"`
	Gender              int8   `json:"gender"`
	Phone               string `json:"phone"`
	Birthday            string `json:"birthday"`
	Age                 int    `json:"age"`
	Zodiac              string `json:"zodiac"`
	Height              int    `json:"height"`
	Weight              int    `json:"weight"`
	Education           int8   `json:"education"`
	MaritalStatus       int8   `json:"marital_status"`
	Profession          string `json:"profession"`
	WorkUnit            string `json:"work_unit"`
	WorkCity            string `json:"work_city"`
	Position            string `json:"position"`
	Income              int    `json:"income"`
	HouseStatus         int8   `json:"house_status"`
	HouseAddress        string `json:"house_address"`
	CarStatus           int8   `json:"car_status"`
	Address             string `json:"address"`
	FamilyDescription   string `json:"family_description"`
	PartnerRequirements string `json:"partner_requirements"`
	ParentsProfession   string `json:"parents_profession"`
	Remark              string `json:"remark"`
	RawText             string `json:"raw_text"`
	ParseStatus         string `json:"parse_status"`
	ErrorMsg            string `json:"error_msg"`
}

type ChatParser struct {
	Records  []ImportRecord
	provider LLMProvider
}

type LLMProvider interface {
	Parse(content string, prompt string) ([]ImportRecord, error)
	Name() string
}

func NewChatParser() *ChatParser {
	provider := getLLMProvider()
	return &ChatParser{
		Records:  make([]ImportRecord, 0),
		provider: provider,
	}
}

func getLLMProvider() LLMProvider {
	cfg := conf.GetConfig()
	if cfg == nil || cfg.LLM == nil {
		log.Warn("LLM config not found, using mock provider")
		return &MockProvider{}
	}

	provider := cfg.LLM.Provider
	if provider == "" {
		provider = "volcano"
	}

	if provider == "volcano" {
		if cfg.LLM.VolcanoEngine != nil && cfg.LLM.VolcanoEngine.APIKey != "" {
			log.Infof("Using Volcano Engine provider: %s", cfg.LLM.VolcanoEngine.Model)
			return &VolcanoProvider{
				APIKey:   cfg.LLM.VolcanoEngine.APIKey,
				Model:    cfg.LLM.VolcanoEngine.Model,
				Endpoint: cfg.LLM.VolcanoEngine.Endpoint,
			}
		}
	}

	log.Warn("No valid LLM provider found, using mock provider")
	return &MockProvider{}
}

type VolcanoProvider struct {
	APIKey   string
	Model    string
	Endpoint string
}

func (p *VolcanoProvider) Name() string {
	return "volcano"
}

func (p *VolcanoProvider) Parse(content string, prompt string) ([]ImportRecord, error) {
	fullPrompt := fmt.Sprintf(prompt, content)

	endpoint := p.Endpoint
	if endpoint == "" {
		endpoint = "https://ark.cn-beijing.volces.com/api/coding/v3/chat/completions"
	}

	reqBody := map[string]interface{}{
		"model": p.Model,
		"messages": []map[string]string{
			{"role": "user", "content": fullPrompt},
		},
	}

	jsonBody, _ := json.Marshal(reqBody)

	client := &http.Client{Timeout: 180 * time.Second}
	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.APIKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Volcano Engine API request failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var volcanoResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&volcanoResp); err != nil {
		log.Errorf("Volcano Engine API response decode failed: %v", err)
		return nil, err
	}

	choices, ok := volcanoResp["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return nil, fmt.Errorf("empty response from Volcano Engine")
	}

	firstChoice, ok := choices[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid choice format from Volcano Engine")
	}

	message, ok := firstChoice["message"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid message format from Volcano Engine")
	}

	contentStr, ok := message["content"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid content format from Volcano Engine")
	}

	return parseResponse(contentStr)
}

type MockProvider struct{}

func (p *MockProvider) Name() string {
	return "mock"
}

func (p *MockProvider) Parse(content string, prompt string) ([]ImportRecord, error) {
	log.Warn("Using mock parser - no LLM API key configured")
	return []ImportRecord{
		{
			Name:                "张三(Mock)",
			Gender:              1,
			Phone:               "13800138000",
			Birthday:            "1990-01",
			Age:                 34,
			Zodiac:              "马",
			Height:              175,
			Weight:              70,
			Education:           3,
			MaritalStatus:       1,
			Profession:          "软件工程师",
			Income:              20000,
			HouseStatus:         2,
			HouseAddress:        "北京市朝阳区",
			CarStatus:           2,
			Address:             "北京市海淀区",
			FamilyDescription:   "父母健在，独生子",
			PartnerRequirements: "本科以上学历，身高160+",
			ParentsProfession:   "父亲教师，母亲家庭主妇",
			Remark:              "不抽烟不喝酒",
			ParseStatus:         "success",
			RawText:             "Mock Data",
		},
	}, nil
}

func parseResponse(contentStr string) ([]ImportRecord, error) {
	contentStr = strings.TrimPrefix(contentStr, "```json")
	contentStr = strings.TrimPrefix(contentStr, "```")
	contentStr = strings.TrimSuffix(contentStr, "```")
	contentStr = strings.TrimSpace(contentStr)

	var records []ImportRecord
	if err := json.Unmarshal([]byte(contentStr), &records); err != nil {
		log.Errorf("JSON parse failed: %v, content: %s", err, contentStr)
		return nil, fmt.Errorf("failed to parse JSON from LLM")
	}

	for i := range records {
		records[i].ParseStatus = "success"
		if records[i].Name == "" {
			records[i].Name = "未知用户"
			records[i].ParseStatus = "warning"
			records[i].ErrorMsg = "未识别到姓名"
		}
		if records[i].Age == 0 && records[i].Birthday != "" && len(records[i].Birthday) >= 7 {
			if age := calculateAge(records[i].Birthday); age > 0 {
				records[i].Age = age
			}
		}
	}

	return records, nil
}

func (p *ChatParser) Parse(content string) ([]ImportRecord, error) {
	prompt := buildParsePrompt()
	return p.provider.Parse(content, prompt)
}

func calculateAge(birthday string) int {
	if len(birthday) < 7 {
		return 0
	}
	birthYear := 0
	fmt.Sscanf(birthday[:4], "%d", &birthYear)
	if birthYear == 0 {
		return 0
	}
	currentYear := time.Now().Year()
	return currentYear - birthYear
}

func buildParsePrompt() string {
	return `【角色】
你是一位专业的婚恋信息录入助手，擅长从聊天记录中提取客户详细信息。

【任务】
请从以下聊天记录中提取所有客户信息，并输出为JSON数组格式。
一条聊天记录可能包含多个客户，请分别提取。

【字段映射规则】
请识别以下所有字段（按优先级排序）：

1. 基础信息：
   - name (string): 姓名，识别关键词：姓名、名字、我叫、我是
   - gender (int): 性别，1=男, 2=女。识别：男/女、先生/女士、帅哥/美女
   - phone (string): 手机号，11位数字
   - birthday (string): 出生年月，格式YYYY-MM。识别：生日、出生日期
   - age (int): 年龄，数字。识别：X岁、年龄
   - zodiac (string): 属相，如：鼠、牛、虎、兔、龙、蛇、马、羊、猴、鸡、狗、猪

2. 身体特征：
   - height (int): 身高cm。识别：身高、Xcm、X厘米
   - weight (int): 体重kg。识别：体重、Xkg、X公斤

3. 教育婚姻：
   - education (int): 学历编码：1高中, 2大专, 3本科, 4硕士, 5博士
   - marital_status (int): 婚姻状况：1未婚, 2已婚, 3离异, 4丧偶

4. 工作收入：
   - profession (string): 职业/工作。识别：职业、工作、职位、干什么、单位、具体工作、做什么工作
   - work_city (string): 工作城市/具体工作地点。识别：具体地点、工作地点、在哪上班、工作城市
   - income (int): 月收入(元)。识别：收入、月薪、工资、X万/月、Xk、税后收入、税后月薪

5. 房产车辆：
   - house_status (int): 房产：1无房, 2已购房, 3贷款购房
   - house_address (string): 买房地址/房产位置/房在哪。识别：房在哪、房子在哪、房产位置、买房地址、房屋地址
   - car_status (int): 车辆：1无车, 2有车。识别：是否有车、车、车情况

6. 家庭信息：
   - address (string): 家庭住址/现居地
   - family_description (string): 家庭成员描述（父母、兄弟姐妹情况）

7. 择偶要求：
   - partner_requirements (string): 对另一半的要求（年龄、身高、学历、地域等）
   - parents_profession (string): 父母工作/职业情况

8. 其他备注：
   - remark (string): 其他需要记录的信息

【重要规则】
1. 如果某个字段在文本中没有提及，设置为null或空字符串
2. 性别、学历、婚姻状况、房产、车辆必须输出数字编码，不要输出文本
3. 收入如果是"1万"，转换为10000；"15k"转换为15000
4. 身高如果是"1.75米"，转换为175
5. 如果一条记录包含多个客户，请拆分成多个JSON对象
6. 只输出JSON数组，不要任何解释文字

【输出示例】
[
  {
    "name": "张三",
    "gender": 1,
    "phone": "13800138000",
    "birthday": "1990-05",
    "age": 34,
    "zodiac": "马",
    "height": 175,
    "weight": 70,
    "education": 3,
    "marital_status": 1,
    "profession": "软件工程师",
    "income": 20000,
    "house_status": 2,
    "house_address": "北京市朝阳区",
    "car_status": 2,
    "address": "北京市海淀区",
    "family_description": "父母健在，独生子",
    "partner_requirements": "希望找本科以上学历，身高160以上的女生",
    "parents_profession": "父亲退休教师，母亲家庭主妇",
    "remark": "不抽烟不喝酒"
  }
]

【待处理文本】
"""
%s
"""`
}

// llm:
//   provider: "volcano"  # volcano
//   volcano_engine:
//     api_key: "5e47477b-73ed-43b4-b9a0-c58f1bbecab2"
//     model: "ep-20260327175716-s5bz9"
//     endpoint: "https://ark.cn-beijing.volces.com/api/v3/chat/completions"
