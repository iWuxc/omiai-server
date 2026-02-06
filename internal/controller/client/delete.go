package client

import (
	"strconv"

	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/pkg/response"

	"github.com/gin-gonic/gin"
)

// Delete 删除客户
func (c *Controller) Delete(ctx *gin.Context) {
	// 从URL参数获取ID
	idStr := ctx.Param("id")
	if idStr == "" {
		response.ErrorResponse(ctx, response.ParamsCommonError, "缺少客户ID")
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.ErrorResponse(ctx, response.ParamsCommonError, "客户ID格式错误")
		return
	}

	// 检查客户是否存在
	client, err := c.client.Get(ctx, id)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "客户不存在")
		return
	}

	// 检查客户是否处于匹配中或已匹配状态
	if client.Status == biz_omiai.ClientStatusMatching || client.Status == biz_omiai.ClientStatusMatched {
		response.ErrorResponse(ctx, response.ParamsCommonError, "该客户正处于匹配状态，无法删除")
		return
	}

	// 执行删除
	if err := c.client.Delete(ctx, id); err != nil {
		response.ErrorResponse(ctx, response.DBDeleteCommonError, "删除客户失败")
		return
	}

	response.SuccessResponse(ctx, "删除成功", nil)
}
