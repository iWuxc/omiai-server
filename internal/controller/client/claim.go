package client

import (
	"omiai-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type ClaimRequest struct {
	ClientID uint64 `json:"client_id" binding:"required"`
}

// Claim 认领客户到私有库
func (c *Controller) Claim(ctx *gin.Context) {
	var req ClaimRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateError(ctx, err, response.ValidateCommonError)
		return
	}

	currentUserID := ctx.GetUint64("current_user_id")
	if currentUserID == 0 {
		response.ErrorResponse(ctx, response.FuncCommonError, "未登录或会话已失效")
		return
	}

	client, err := c.Client.Get(ctx, req.ClientID)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "客户不存在")
		return
	}

	if !client.IsPublic {
		response.ErrorResponse(ctx, response.ParamsCommonError, "该客户已被认领")
		return
	}

	// 执行认领：更新 manager_id 并设为非公海
	client.ManagerID = currentUserID
	client.IsPublic = false
	if err := c.Client.Update(ctx, client); err != nil {
		response.ErrorResponse(ctx, response.DBUpdateCommonError, "认领失败")
		return
	}

	response.SuccessResponse(ctx, "认领成功", nil)
}

// Release 释放客户到公海池
func (c *Controller) Release(ctx *gin.Context) {
	var req ClaimRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateError(ctx, err, response.ValidateCommonError)
		return
	}

	currentUserID := ctx.GetUint64("current_user_id")

	client, err := c.Client.Get(ctx, req.ClientID)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "客户不存在")
		return
	}

	// 权限校验：只能释放自己的客户，或者是管理员(假设 admin role=1)
	currentUserRole := ctx.GetInt("current_user_role")
	if client.ManagerID != currentUserID && currentUserRole != 1 {
		response.ErrorResponse(ctx, response.AuthCommonError, "无权操作此客户")
		return
	}

	// 执行释放：manager_id=0, is_public=true
	client.ManagerID = 0
	client.IsPublic = true
	if err := c.Client.Update(ctx, client); err != nil {
		response.ErrorResponse(ctx, response.DBUpdateCommonError, "释放失败")
		return
	}

	response.SuccessResponse(ctx, "释放成功", nil)
}
