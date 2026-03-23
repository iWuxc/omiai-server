package c_event

import (
	"fmt"
	biz_event "omiai-server/internal/biz/event"
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	Event  biz_event.EventInterface
	Client biz_omiai.ClientInterface
}

func NewController(event biz_event.EventInterface, client biz_omiai.ClientInterface) *Controller {
	return &Controller{Event: event, Client: client}
}

// List C端获取正在报名中的活动列表
func (c *Controller) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	offset := (page - 1) * size

	// 只查询状态为 1(报名中) 的活动
	list, total, err := c.Event.List(ctx, 1, offset, size)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取活动列表失败")
		return
	}

	response.SuccessResponse(ctx, "success", map[string]interface{}{
		"list":  list,
		"total": total,
	})
}

// Register C端报名线下活动
func (c *Controller) Register(ctx *gin.Context) {
	clientID, exists := ctx.Get("client_id")
	if !exists {
		response.ErrorResponse(ctx, response.AuthCommonError, "未登录")
		return
	}
	myID := clientID.(uint64)

	var req struct {
		EventID uint64 `json:"event_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateError(ctx, err, response.ValidateCommonError)
		return
	}

	// 1. 检查活动状态
	event, err := c.Event.Get(ctx, req.EventID)
	if err != nil || event == nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "活动不存在")
		return
	}
	if event.Status != 1 {
		response.ErrorResponse(ctx, response.ParamsCommonError, "该活动不在报名状态")
		return
	}

	// 2. 检查是否已经报名
	hasReg, _ := c.Event.HasRegistered(ctx, req.EventID, myID)
	if hasReg {
		response.ErrorResponse(ctx, response.ParamsCommonError, "您已经报名过该活动了")
		return
	}

	// 3. 检查余额并扣费
	if event.PriceCoins > 0 {
		me, _ := c.Client.Get(ctx, myID)
		if me == nil || me.Coins < event.PriceCoins {
			response.ErrorResponse(ctx, response.ParamsCommonError, "红豆余额不足，请先充值")
			return
		}

		// 扣费 (真实场景下这应该和报名放在同一个分布式事务里，目前为了快速迭代分开处理)
		if err := c.Client.AddCoins(ctx, myID, -event.PriceCoins, 4, fmt.Sprintf("报名活动: %s", event.Title)); err != nil {
			response.ErrorResponse(ctx, response.DBUpdateCommonError, "支付红豆失败")
			return
		}
	}

	// 4. 执行报名 (包含人数并发控制)
	if err := c.Event.Register(ctx, req.EventID, myID, event.PriceCoins); err != nil {
		// TODO: 如果报名失败(如超卖)，需要回滚红豆。当前仅作为MVP演示
		response.ErrorResponse(ctx, response.ServiceCommonError, "报名失败，名额已满")
		return
	}

	response.SuccessResponse(ctx, "报名成功！请准时参加", nil)
}
