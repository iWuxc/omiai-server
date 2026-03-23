package c_pay

import (
	"fmt"
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/data"
	"omiai-server/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	db     *data.DB
	Client biz_omiai.ClientInterface
}

func NewController(db *data.DB, client biz_omiai.ClientInterface) *Controller {
	return &Controller{
		db:     db,
		Client: client,
	}
}

type RechargeRequest struct {
	Amount int `json:"amount" binding:"required,min=1"` // 充值金额，对应红豆数量，例如1元=10红豆
}

// Recharge 模拟充值红豆 (实际生产需对接微信支付并在回调中处理)
func (c *Controller) Recharge(ctx *gin.Context) {
	clientID, exists := ctx.Get("client_id")
	if !exists {
		response.ErrorResponse(ctx, response.AuthCommonError, "未授权")
		return
	}

	var req RechargeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateError(ctx, err, response.ValidateCommonError)
		return
	}

	coinsToAdd := req.Amount * 10 // 充值比例 1:10

	if err := c.Client.AddCoins(ctx, clientID.(uint64), coinsToAdd, 1, fmt.Sprintf("充值%d元", req.Amount)); err != nil {
		response.ErrorResponse(ctx, response.DBUpdateCommonError, "充值失败")
		return
	}

	response.SuccessResponse(ctx, "充值成功", map[string]int{"coins_added": coinsToAdd})
}

type UnlockRequest struct {
	TargetClientID uint64 `json:"target_client_id" binding:"required"`
}

// UnlockProfile 使用红豆解锁查看喜欢我的人
func (c *Controller) UnlockProfile(ctx *gin.Context) {
	clientID, exists := ctx.Get("client_id")
	if !exists {
		response.ErrorResponse(ctx, response.AuthCommonError, "未授权")
		return
	}
	myID := clientID.(uint64)

	var req UnlockRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateError(ctx, err, response.ValidateCommonError)
		return
	}

	// 如果是 VIP，免费解锁 (此处为了演示简单，只要是VIP就不扣费)
	isVip := c.Client.IsVip(ctx, myID)

	unlockCost := 20 // 每次解锁消耗 20 红豆 (约合2元)
	if !isVip {
		me, _ := c.Client.Get(ctx, myID)
		if me == nil || me.Coins < unlockCost {
			response.ErrorResponse(ctx, response.ParamsCommonError, "红豆余额不足，请先充值")
			return
		}

		// 扣费
		if err := c.Client.AddCoins(ctx, myID, -unlockCost, 2, fmt.Sprintf("解锁查看用户 %d", req.TargetClientID)); err != nil {
			response.ErrorResponse(ctx, response.DBUpdateCommonError, "扣费失败")
			return
		}
	}

	// 扣费成功，返回对方详细信息
	targetClient, err := c.Client.Get(ctx, req.TargetClientID)
	if err != nil || targetClient == nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "目标用户不存在")
		return
	}

	var cost interface{}
	if isVip {
		cost = 0
	} else {
		cost = unlockCost
	}

	response.SuccessResponse(ctx, "解锁成功", map[string]interface{}{
		"id":         targetClient.ID,
		"name":       targetClient.Name,
		"avatar":     targetClient.Avatar,
		"phone":      targetClient.Phone,
		"wechat":     targetClient.WxOpenid,
		"cost_coins": cost,
	})
}

// BuyVip 购买VIP特权
func (c *Controller) BuyVip(ctx *gin.Context) {
	clientID, exists := ctx.Get("client_id")
	if !exists {
		response.ErrorResponse(ctx, response.AuthCommonError, "未授权")
		return
	}
	myID := clientID.(uint64)

	vipCost := 990 // 包月VIP 99元 = 990红豆
	me, _ := c.Client.Get(ctx, myID)
	if me == nil || me.Coins < vipCost {
		response.ErrorResponse(ctx, response.ParamsCommonError, "红豆余额不足，请先充值")
		return
	}

	// 扣费
	if err := c.Client.AddCoins(ctx, myID, -vipCost, 3, "开通包月VIP"); err != nil {
		response.ErrorResponse(ctx, response.DBUpdateCommonError, "扣费失败")
		return
	}

	// 更新 VIP 到期时间
	newExpireAt := time.Now().AddDate(0, 1, 0)
	if me.VipExpireAt.After(time.Now()) {
		newExpireAt = me.VipExpireAt.AddDate(0, 1, 0)
	}
	me.VipExpireAt = newExpireAt

	if err := c.Client.Update(ctx, me); err != nil {
		response.ErrorResponse(ctx, response.DBUpdateCommonError, "更新VIP状态失败")
		return
	}

	response.SuccessResponse(ctx, "开通VIP成功", map[string]interface{}{
		"vip_expire_at": me.VipExpireAt.Format("2006-01-02 15:04:05"),
	})
}
