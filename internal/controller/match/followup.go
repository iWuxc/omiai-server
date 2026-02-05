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

	var followUpDate time.Time
	if req.FollowUpDate == "" {
		followUpDate = time.Now()
	} else {
		var err error
		followUpDate, err = parseTime(req.FollowUpDate)
		if err != nil {
			response.ErrorResponse(ctx, response.ParamsCommonError, "回访时间格式错误: "+err.Error())
			return
		}
	}

	var nextFollowUpAt time.Time
	if req.NextFollowUpAt != "" {
		var err error
		nextFollowUpAt, err = parseTime(req.NextFollowUpAt)
		if err != nil {
			response.ErrorResponse(ctx, response.ParamsCommonError, "下次回访时间格式错误: "+err.Error())
			return
		}
	}

	record := &biz_omiai.FollowUpRecord{
		MatchRecordID:  req.MatchRecordID,
		FollowUpDate:   followUpDate,
		Method:         req.Method,
		Content:        req.Content,
		Feedback:       req.Feedback,
		Satisfaction:   req.Satisfaction,
		Attachments:    req.Attachments,
		NextFollowUpAt: nextFollowUpAt,
	}

	if err := c.match.CreateFollowUp(ctx, record); err != nil {
		response.ErrorResponse(ctx, response.DBInsertCommonError, "保存回访记录失败")
		return
	}

	response.SuccessResponse(ctx, "保存成功", record)
}

func parseTime(val string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, val); err == nil {
		return t, nil
	}
	if t, err := time.ParseInLocation("2006-01-02 15:04:05", val, time.Local); err == nil {
		return t, nil
	}
	if t, err := time.ParseInLocation("2006-01-02", val, time.Local); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("supported formats: RFC3339, YYYY-MM-DD HH:mm:ss, YYYY-MM-DD")
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
