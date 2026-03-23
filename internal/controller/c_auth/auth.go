package c_auth

import (
	"fmt"
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/data"
	"omiai-server/pkg/auth"
	"omiai-server/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iWuxc/go-wit/log"
	"github.com/iWuxc/go-wit/redis"
)

type Controller struct {
	db     *data.DB
	Client biz_omiai.ClientInterface
	Redis  *redis.Redis
}

func NewController(db *data.DB, client biz_omiai.ClientInterface) *Controller {
	return &Controller{
		db:     db,
		Client: client,
		Redis:  redis.GetRedis(),
	}
}

type WxLoginRequest struct {
	Code string `json:"code" binding:"required"`
}

// WxLogin C端微信登录
func (c *Controller) WxLogin(ctx *gin.Context) {
	var req WxLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateError(ctx, err, response.ValidateCommonError)
		return
	}

	// 1. 调用微信接口换取 openid (此处使用 Mock)
	// TODO: 替换为真实的微信 API 调用
	openID := mockWechatLogin(req.Code)
	if openID == "" {
		response.ErrorResponse(ctx, response.ParamsCommonError, "微信登录失败")
		return
	}

	// 2. 根据 openid 查找用户
	client, err := c.Client.GetByWxOpenID(ctx, openID)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "系统错误")
		return
	}

	isNew := false
	if client == nil {
		// 3. 如果用户不存在，静默注册
		isNew = true
		client = &biz_omiai.Client{
			WxOpenid:   openID,
			Name:       "微信用户_" + time.Now().Format("0102150405"),
			Status:     biz_omiai.ClientStatusSingle,
			IsPublic:   true,
			IsVerified: false,
		}
		if err := c.Client.Create(ctx, client); err != nil {
			log.Errorf("Create CClient failed: %v", err)
			response.ErrorResponse(ctx, response.DBInsertCommonError, "注册失败")
			return
		}
	}

	// 4. 生成 Token
	token, err := auth.GenerateToken(client.ID, "c_client")
	if err != nil {
		response.ErrorResponse(ctx, response.ServiceCommonError, "生成Token失败")
		return
	}

	response.SuccessResponse(ctx, "登录成功", map[string]interface{}{
		"token":  token,
		"is_new": isNew,
		"client": client,
	})
}

func mockWechatLogin(code string) string {
	// 在没有真实微信 AppID 和 Secret 之前，用 code 直接生成一个 mock openid
	return fmt.Sprintf("mock_openid_%s", code)
}
