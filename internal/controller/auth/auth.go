package auth

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
	db    *data.DB
	User  biz_omiai.UserInterface
	Redis *redis.Redis
}

func NewController(db *data.DB, user biz_omiai.UserInterface) *Controller {
	return &Controller{
		db:    db,
		User:  user,
		Redis: redis.GetRedis(),
	}
}

type SmsLoginRequest struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

type SendSmsRequest struct {
	Phone string `json:"phone" binding:"required"`
}

// SendSms 发送验证码
func (c *Controller) SendSms(ctx *gin.Context) {
	var req SendSmsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateError(ctx, err, response.ValidateCommonError)
		return
	}

	// 1. 防刷：限制发送频率 (e.g., 1 minute)
	lockKey := fmt.Sprintf("sms:lock:%s", req.Phone)
	if c.Redis.GetClient().Exists(ctx, lockKey).Val() > 0 {
		response.ErrorResponse(ctx, response.ParamsCommonError, "发送过于频繁，请稍后再试")
		return
	}

	// 2. 生成验证码
	code := "123456" // Mock code for now

	// 3. 存储到 Redis (有效期 5 分钟)
	codeKey := fmt.Sprintf("sms:code:%s", req.Phone)
	c.Redis.GetClient().Set(ctx, codeKey, code, 5*time.Minute)
	c.Redis.GetClient().Set(ctx, lockKey, 1, 1*time.Minute)

	log.Infof("Sent SMS code %s to %s", code, req.Phone)
	response.SuccessResponse(ctx, "验证码已发送", nil)
}

// H5Login H5 登录 (手机号 + 验证码)
func (c *Controller) H5Login(ctx *gin.Context) {
	var req SmsLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateError(ctx, err, response.ValidateCommonError)
		return
	}

	// 1. 验证码校验
	// 特殊处理测试账号
	if !(req.Phone == "13800138000" && req.Code == "123456") {
		codeKey := fmt.Sprintf("sms:code:%s", req.Phone)
		savedCode := c.Redis.GetClient().Get(ctx, codeKey).Val()
		if savedCode == "" || savedCode != req.Code {
			response.ErrorResponse(ctx, response.ParamsCommonError, "验证码错误或已过期")
			return
		}
		// 校验成功后清除验证码
		c.Redis.GetClient().Del(ctx, codeKey)
	}

	// 2. 获取或创建用户
	user, err := c.User.GetByPhone(ctx, req.Phone)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "系统错误")
		return
	}

	if user == nil {
		user = &biz_omiai.User{
			Phone:    req.Phone,
			Nickname: fmt.Sprintf("用户_%s", req.Phone[len(req.Phone)-4:]),
			Role:     biz_omiai.RoleOperator,
		}
		if err := c.User.Create(ctx, user); err != nil {
			response.ErrorResponse(ctx, response.DBInsertCommonError, "创建用户失败")
			return
		}
	}

	// 3. 生成 Token
	token, err := auth.GenerateToken(user.ID, user.Role)
	if err != nil {
		response.ErrorResponse(ctx, response.FuncCommonError, "生成 Token 失败")
		return
	}

	response.SuccessResponse(ctx, "登录成功", map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

type WxLoginRequest struct {
	Code string `json:"code" binding:"required"`
}

// WxLogin 小程序登录
func (c *Controller) WxLogin(ctx *gin.Context) {
	var req WxLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateError(ctx, err, response.ValidateCommonError)
		return
	}

	// 1. 调用微信接口获取 OpenID (Mock for now)
	openID := fmt.Sprintf("mock_openid_%s", req.Code)

	// 2. 获取或创建用户
	user, err := c.User.GetByWxOpenID(ctx, openID)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "系统错误")
		return
	}

	if user == nil {
		user = &biz_omiai.User{
			WxOpenID: openID,
			Nickname: "微信用户",
			Role:     biz_omiai.RoleOperator,
		}
		if err := c.User.Create(ctx, user); err != nil {
			response.ErrorResponse(ctx, response.DBInsertCommonError, "创建用户失败")
			return
		}
	}

	// 3. 生成 Token
	token, err := auth.GenerateToken(user.ID, user.Role)
	if err != nil {
		response.ErrorResponse(ctx, response.FuncCommonError, "生成 Token 失败")
		return
	}

	response.SuccessResponse(ctx, "登录成功", map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

// GetUserInfo 获取当前用户信息
func (c *Controller) GetUserInfo(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	user, err := c.User.GetByID(ctx, userID)
	if err != nil || user == nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "用户不存在")
		return
	}

	response.SuccessResponse(ctx, "ok", user)
}
