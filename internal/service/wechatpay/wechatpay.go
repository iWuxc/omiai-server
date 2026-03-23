package wechatpay

import (
	"context"
	"fmt"
	"time"

	"github.com/iWuxc/go-wit/log"
)

// PayOrder 支付订单参数
type PayOrder struct {
	OutTradeNo  string // 商户订单号
	Description string // 商品描述
	Amount      int    // 总金额，单位：分
	OpenID      string // 用户OpenID
}

// PayResult 支付结果
type PayResult struct {
	PrepayID string
	NonceStr string
	TimeStamp string
	Sign     string
	SignType string
	AppID    string
}

// Service 微信支付服务接口
type Service interface {
	// CreateMiniProgramOrder 创建小程序支付订单
	CreateMiniProgramOrder(ctx context.Context, order *PayOrder) (*PayResult, error)
	// VerifyCallback 验证支付回调签名并解析
	VerifyCallback(ctx context.Context, reqBody []byte, signature, timestamp, nonce string) (map[string]interface{}, error)
}

type WechatPayService struct {
	AppID     string
	MchID     string
	APIv3Key  string
	CertPath  string
	NotifyURL string
}

func NewWechatPayService() Service {
	// TODO: 生产环境应从配置中心加载这些敏感信息，并使用官方 SDK "github.com/wechatpay-apiv3/wechatpay-go"
	return &WechatPayService{
		AppID:     "wx_mock_appid_12345",
		MchID:     "1234567890",
		APIv3Key:  "mock_apiv3_key",
		NotifyURL: "https://api.omiai.com/api/c/pay/wechat_notify",
	}
}

func (s *WechatPayService) CreateMiniProgramOrder(ctx context.Context, order *PayOrder) (*PayResult, error) {
	log.Infof("[WechatPay] Creating Order: %s, Amount: %d, OpenID: %s", order.OutTradeNo, order.Amount, order.OpenID)

	// 模拟调用微信支付 JSAPI 下单接口
	// 实际生产： client.Post(ctx, "https://api.mch.weixin.qq.com/v3/pay/transactions/jsapi", ...)
	
	// 模拟返回签名供小程序拉起收银台
	return &PayResult{
		AppID:     s.AppID,
		PrepayID:  fmt.Sprintf("wx%d", time.Now().UnixNano()),
		NonceStr:  "mock_nonce_str",
		TimeStamp: fmt.Sprintf("%d", time.Now().Unix()),
		SignType:  "RSA",
		Sign:      "mock_rsa_signature",
	}, nil
}

func (s *WechatPayService) VerifyCallback(ctx context.Context, reqBody []byte, signature, timestamp, nonce string) (map[string]interface{}, error) {
	log.Infof("[WechatPay] Verifying Callback Signature...")
	
	// 模拟验证通过并解析解密后的报文
	// 实际生产：使用官方 SDK 的 core.VerifySign(...) 和 utils.DecryptAES256GCM(...)
	
	mockDecryptedData := map[string]interface{}{
		"out_trade_no":   "PAY_20260323_123456",
		"transaction_id": "4200000000000000000000000000",
		"trade_state":    "SUCCESS",
		"amount": map[string]interface{}{
			"total": 990, // 9.9元
		},
	}
	return mockDecryptedData, nil
}
