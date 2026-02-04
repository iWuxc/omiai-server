package match

import (
	"fmt"
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/validates"
	"omiai-server/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
)

func (c *Controller) CreateFollowUp(ctx *gin.Context) {
	var req validates.FollowUpCreateValidate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateError(ctx, err, response.ValidateCommonError)
		return
	}

	followUpDate := req.FollowUpDate
	if followUpDate.IsZero() {
		followUpDate = time.Now()
	}

	record := &biz_omiai.FollowUpRecord{
		MatchRecordID:  req.MatchRecordID,
		FollowUpDate:   followUpDate,
		Method:         req.Method,
		Content:        req.Content,
		Feedback:       req.Feedback,
		Satisfaction:   req.Satisfaction,
		Attachments:    req.Attachments,
		NextFollowUpAt: req.NextFollowUpAt,
	}

	if err := c.match.CreateFollowUp(ctx, record); err != nil {
		response.ErrorResponse(ctx, response.DBInsertCommonError, "保存回访记录失败")
		return
	}

	response.SuccessResponse(ctx, "保存成功", record)
}

func (c *Controller) ListFollowUps(ctx *gin.Context) {
	idStr := ctx.Query("match_record_id")
	if idStr == "" {
		response.ErrorResponse(ctx, response.ParamsCommonError, "参数错误")
		return
	}

	// In a real app, use a proper parser
	var id uint64
	fmt.Sscanf(idStr, "%d", &id)

	list, err := c.match.SelectFollowUps(ctx, id)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "查询失败")
		return
	}

	response.SuccessResponse(ctx, "ok", list)
}

func (c *Controller) GetReminders(ctx *gin.Context) {
	list, err := c.match.GetReminders(ctx)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "查询失败")
		return
	}
	response.SuccessResponse(ctx, "ok", list)
}
