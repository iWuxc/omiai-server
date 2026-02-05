package match

import (
	"fmt"
	"omiai-server/internal/validates"
	"omiai-server/pkg/response"

	"github.com/gin-gonic/gin"
)

// Confirm 确认匹配 (管理员直接确认)
func (c *Controller) Confirm(ctx *gin.Context) {
	var req validates.ConfirmMatchValidate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateError(ctx, err, response.ValidateCommonError)
		return
	}

	// Get admin ID from context (set by middleware)
	adminID := "unknown"
	if v, exists := ctx.Get("user_id"); exists {
		adminID = fmt.Sprintf("%v", v)
	}

	matchRecord, err := c.match.ConfirmMatch(ctx, req.ClientID, req.CandidateID, adminID)
	if err != nil {
		response.ErrorResponse(ctx, response.DBInsertCommonError, "确认匹配失败")
		return
	}

	response.SuccessResponse(ctx, "匹配确认成功", matchRecord)
}
