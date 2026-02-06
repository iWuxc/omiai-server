package ai

import (
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/data"
	aiservice "omiai-server/internal/service/ai"
	"omiai-server/internal/validates"
	"omiai-server/pkg/response"

	"github.com/gin-gonic/gin"
)

// Controller AI分析控制器
type Controller struct {
	db         *data.DB
	clientRepo biz_omiai.ClientInterface
	aiAnalyzer *aiservice.AIAnalyzer
}

// NewController 创建AI控制器
func NewController(db *data.DB, clientRepo biz_omiai.ClientInterface) *Controller {
	return &Controller{
		db:         db,
		clientRepo: clientRepo,
		aiAnalyzer: aiservice.NewAIAnalyzer(),
	}
}

// AnalyzeMatch AI匹配分析
func (c *Controller) AnalyzeMatch(ctx *gin.Context) {
	var req validates.AIAnalyzeValidate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateError(ctx, err, response.ValidateCommonError)
		return
	}

	// 获取两个客户信息
	clientA, err := c.clientRepo.Get(ctx, req.ClientAID)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取客户A信息失败")
		return
	}

	clientB, err := c.clientRepo.Get(ctx, req.ClientBID)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取客户B信息失败")
		return
	}

	// 转换为AI分析用的格式
	profileA := convertToProfile(clientA)
	profileB := convertToProfile(clientB)

	// 调用AI分析
	result, err := c.aiAnalyzer.AnalyzeMatch(profileA, profileB)
	if err != nil {
		response.ErrorResponse(ctx, response.FuncCommonError, "AI分析失败："+err.Error())
		return
	}

	response.SuccessResponse(ctx, "分析完成", result)
}

// GetIceBreaker 获取破冰话题
func (c *Controller) GetIceBreaker(ctx *gin.Context) {
	var req validates.AIAnalyzeValidate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateError(ctx, err, response.ValidateCommonError)
		return
	}

	clientA, err := c.clientRepo.Get(ctx, req.ClientAID)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取客户A信息失败")
		return
	}

	clientB, err := c.clientRepo.Get(ctx, req.ClientBID)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取客户B信息失败")
		return
	}

	profileA := convertToProfile(clientA)
	profileB := convertToProfile(clientB)

	topics, err := c.aiAnalyzer.GenerateIceBreaker(profileA, profileB)
	if err != nil {
		response.ErrorResponse(ctx, response.FuncCommonError, "生成话题失败")
		return
	}

	response.SuccessResponse(ctx, "获取成功", gin.H{
		"topics": topics,
	})
}

// convertToProfile 将Client转换为AI分析用的Profile
func convertToProfile(client *biz_omiai.Client) *aiservice.ClientProfile {
	gender := "男"
	if client.Gender == 2 {
		gender = "女"
	}

	education := getEducationText(client.Education)
	marital := getMaritalText(client.MaritalStatus)
	house := getHouseText(client.HouseStatus)
	car := getCarText(client.CarStatus)

	return &aiservice.ClientProfile{
		Name:                client.Name,
		Gender:              gender,
		Age:                 client.RealAge(),
		Height:              client.Height,
		Education:           education,
		Income:              client.Income,
		Profession:          client.Profession,
		MaritalStatus:       marital,
		HouseStatus:         house,
		CarStatus:           car,
		Address:             client.Address,
		FamilyDescription:   client.FamilyDescription,
		PartnerRequirements: client.PartnerRequirements,
		Remark:              client.Remark,
		Tags:                client.Tags,
	}
}

func getEducationText(edu int8) string {
	switch edu {
	case 1:
		return "高中及以下"
	case 2:
		return "大专"
	case 3:
		return "本科"
	case 4:
		return "硕士"
	case 5:
		return "博士"
	default:
		return "未知"
	}
}

func getMaritalText(status int8) string {
	switch status {
	case 1:
		return "未婚"
	case 2:
		return "已婚"
	case 3:
		return "离异"
	case 4:
		return "丧偶"
	default:
		return "未知"
	}
}

func getHouseText(status int8) string {
	switch status {
	case 1:
		return "无房"
	case 2:
		return "已购房"
	case 3:
		return "贷款购房"
	default:
		return "未知"
	}
}

func getCarText(status int8) string {
	switch status {
	case 1:
		return "无车"
	case 2:
		return "有车"
	default:
		return "未知"
	}
}
