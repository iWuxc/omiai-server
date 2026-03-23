package c_pay

import (
	"fmt"
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/data"
	"omiai-server/internal/service/wechatpay"
	"omiai-server/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iWuxc/go-wit/log"
)

type Controller struct {
	db        *data.DB
	Client    biz_omiai.ClientInterface
	WechatPay wechatpay.Service
}

func NewController(db *data.DB, client biz_omiai.ClientInterface, wp wechatpay.Service) *Controller {
	return &Controller{
		db:        db,
		Client:    client,
		WechatPay: wp,
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

	me, _ := c.Client.Get(ctx, clientID.(uint64))
	if me == nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "用户不存在")
		return
	}

	// 1. 生成内部业务订单号
	outTradeNo := fmt.Sprintf("PAY_%d_%d", me.ID, time.Now().UnixNano())

	// 2. 调用微信支付统一下单 (金额单位:分)
	orderAmountFen := req.Amount * 100
	payResult, err := c.WechatPay.CreateMiniProgramOrder(ctx, &wechatpay.PayOrder{
		OutTradeNo:  outTradeNo,
		Description: fmt.Sprintf("充值 %d 红豆", req.Amount*10),
		Amount:      orderAmountFen,
		OpenID:      me.WxOpenid,
	})

	if err != nil {
		log.Errorf("Create wechat order failed: %v", err)
		response.ErrorResponse(ctx, response.ServiceCommonError, "唤起支付失败")
		return
	}

	// TODO: 这里应在数据库 `payment_order` 表记录这条订单，状态为"待支付"

	// 3. 将签名参数返回给小程序端，拉起微信收银台
	response.SuccessResponse(ctx, "下单成功", map[string]interface{}{
		"out_trade_no": outTradeNo,
		"pay_params":   payResult,
	})
}

// WechatNotify 微信支付回调接口 (供微信服务器调用)
func (c *Controller) WechatNotify(ctx *gin.Context) {
	// 1. 读取 Header 里的签名信息
	signature := ctx.GetHeader("Wechatpay-Signature")
	timestamp := ctx.GetHeader("Wechatpay-Timestamp")
	nonce := ctx.GetHeader("Wechatpay-Nonce")

	// 2. 读取 Body
	body, _ := ctx.GetRawData()

	// 3. 验证签名并解密
	payData, err := c.WechatPay.VerifyCallback(ctx, body, signature, timestamp, nonce)
	if err != nil {
		log.Errorf("Wechat notify verify failed: %v", err)
		ctx.JSON(400, gin.H{"code": "FAIL", "message": "Verify failed"})
		return
	}

	// 4. 处理业务逻辑 (发放红豆)
	outTradeNo := payData["out_trade_no"].(string)
	tradeState := payData["trade_state"].(string)

	if tradeState == "SUCCESS" {
		log.Infof("Order %s paid successfully, starting to deliver virtual coins...", outTradeNo)

		// TODO: 从 outTradeNo 解析出 ClientID，或者查询 `payment_order` 表
		// 假设我们解析出 ClientID=1，充值金额 9.9元(990分) -> 发放 99 红豆
		mockClientID := uint64(1)
		mockAmountFen := 990
		coinsToAdd := (mockAmountFen / 100) * 10

		// 防重放：判断该订单是否已经处理过

		// 执行加币
		_ = c.Client.AddCoins(ctx, mockClientID, coinsToAdd, 1, fmt.Sprintf("微信支付充值单号: %s", outTradeNo))
	}

	// 5. 必须返回 200 给微信服务器
	ctx.JSON(200, gin.H{"code": "SUCCESS", "message": "OK"})
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
