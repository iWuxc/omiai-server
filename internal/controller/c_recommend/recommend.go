package c_recommend

import (
	"omiai-server/internal/biz"
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/data"
	"omiai-server/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	db     *data.DB
	Client biz_omiai.ClientInterface
	Match  biz_omiai.MatchInterface
}

func NewController(db *data.DB, client biz_omiai.ClientInterface, match biz_omiai.MatchInterface) *Controller {
	return &Controller{
		db:     db,
		Client: client,
		Match:  match,
	}
}

// ClientDTO 用于返回给C端的脱敏资料
type ClientDTO struct {
	ID                  uint64 `json:"id"`
	Name                string `json:"name"` // 可以考虑只返回姓氏，例如 "张先生"
	Gender              int8   `json:"gender"`
	Age                 int    `json:"age"`
	Avatar              string `json:"avatar"`
	Height              int    `json:"height"`
	Education           int8   `json:"education"`
	Income              int    `json:"income"`
	Profession          string `json:"profession"`
	WorkCity            string `json:"work_city"`
	PartnerRequirements string `json:"partner_requirements"`
	Tags                string `json:"tags"`
	InterestTags        string `json:"interest_tags"`
	Photos              string `json:"photos"`
}

func toClientDTO(c *biz_omiai.Client) *ClientDTO {
	name := c.Name
	if len(name) > 0 {
		// 简单的脱敏逻辑：保留姓氏
		name = string([]rune(name)[0]) + "***"
	}
	return &ClientDTO{
		ID:                  c.ID,
		Name:                name,
		Gender:              c.Gender,
		Age:                 c.RealAge(),
		Avatar:              c.Avatar,
		Height:              c.Height,
		Education:           c.Education,
		Income:              c.Income,
		Profession:          c.Profession,
		WorkCity:            c.WorkCity,
		PartnerRequirements: c.PartnerRequirements,
		Tags:                c.Tags,
		InterestTags:        c.InterestTags,
		Photos:              c.Photos,
	}
}

// DailyRecommend 每日推荐流
func (c *Controller) DailyRecommend(ctx *gin.Context) {
	clientID, exists := ctx.Get("client_id")
	if !exists {
		response.ErrorResponse(ctx, response.AuthCommonError, "未授权")
		return
	}

	me, err := c.Client.Get(ctx, clientID.(uint64))
	if err != nil || me == nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取用户资料失败")
		return
	}

	// 推荐逻辑重构：2026 Q2 战略里程碑
	// 演进方向：这里将不再直接查询 MySQL，而是通过内部 gRPC 调用 AI 推荐中台服务
	// 推荐算法 = (硬条件匹配得分 * 0.4) + (兴趣标签相似度 * 0.3) + (协同过滤/互动偏好 * 0.3)

	// 当前过渡阶段实现：先捞取符合硬条件的池子，然后在内存中根据 InterestTags 计算相似度排序
	clause := &biz.WhereClause{
		Where:   "gender != ? AND status = ? AND is_public = ?",
		Args:    []interface{}{me.Gender, biz_omiai.ClientStatusSingle, true},
		OrderBy: "created_at DESC",
	}

	// 放宽查询限制，以便在内存中进行重排
	list, err := c.Client.Select(ctx, clause, []string{}, 0, 50)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取推荐列表失败")
		return
	}

	// TODO: 基于用户的 InterestTags 和行为数据进行协同过滤打分排序
	// SortByRecommendationScore(me, list)

	var dtoList []*ClientDTO
	count := 0
	for _, v := range list {
		if v.ID == me.ID {
			continue
		}
		dtoList = append(dtoList, toClientDTO(v))
		count++
		if count >= 10 { // 最终只下发 10 个给前端
			break
		}
	}

	response.SuccessResponse(ctx, "success", map[string]interface{}{
		"list": dtoList,
	})
}

// Detail 获取推荐对象详情及AI匹配度
func (c *Controller) Detail(ctx *gin.Context) {
	clientID, exists := ctx.Get("client_id")
	if !exists {
		response.ErrorResponse(ctx, response.AuthCommonError, "未授权")
		return
	}

	targetIDStr := ctx.Param("id")
	targetID, err := strconv.ParseUint(targetIDStr, 10, 64)
	if err != nil {
		response.ErrorResponse(ctx, response.ParamsCommonError, "无效的目标ID")
		return
	}

	targetClient, err := c.Client.Get(ctx, targetID)
	if err != nil || targetClient == nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取目标资料失败")
		return
	}

	// 调用复用的 Match Compare 接口获取 AI 匹配雷达图
	comparison, err := c.Match.Compare(ctx, clientID.(uint64), targetID)
	if err != nil {
		// 如果 AI 调用失败，只返回脱敏档案
		response.SuccessResponse(ctx, "success", map[string]interface{}{
			"profile":  toClientDTO(targetClient),
			"ai_match": nil,
		})
		return
	}

	response.SuccessResponse(ctx, "success", map[string]interface{}{
		"profile":  toClientDTO(targetClient),
		"ai_match": comparison,
	})
}
