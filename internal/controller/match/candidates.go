package match

import (
	"omiai-server/internal/validates"
	"omiai-server/pkg/response"

	"github.com/gin-gonic/gin"
)

// GetCandidates 获取匹配候选人列表
func (c *Controller) GetCandidates(ctx *gin.Context) {
	var req validates.GetCandidatesValidate
	if err := ctx.ShouldBindUri(&req); err != nil {
		response.ValidateError(ctx, err, response.ParamsCommonError)
		return
	}

	// Check if client exists
	client, err := c.client.Get(ctx, req.ClientID)
	if err != nil || client == nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "客户档案不存在")
		return
	}

	candidates, err := c.match.GetCandidates(ctx, req.ClientID)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取候选人失败")
		return
	}

	response.SuccessResponse(ctx, "获取成功", candidates)
}

// Compare 匹配对比详情
func (c *Controller) Compare(ctx *gin.Context) {
	var req validates.CompareValidate
	if err := ctx.ShouldBindUri(&req); err != nil {
		response.ValidateError(ctx, err, response.ParamsCommonError)
		return
	}

	comparison, err := c.match.Compare(ctx, req.ClientID, req.CandidateID)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取对比详情失败")
		return
	}

	response.SuccessResponse(ctx, "获取成功", comparison)
}
