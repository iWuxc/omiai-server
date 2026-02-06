package reminder

import (
	"fmt"
	"omiai-server/internal/validates"
	"omiai-server/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// List 获取提醒列表
func (c *Controller) List(ctx *gin.Context) {
	var req validates.ReminderListValidate
	if err := ctx.ShouldBind(&req); err != nil {
		response.ValidateError(ctx, err, response.ValidateCommonError)
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	offset := (req.Page - 1) * req.PageSize

	// 获取当前用户ID
	userID := ctx.GetUint64("user_id")
	if userID == 0 {
		userID = 1 // 默认用户
	}

	list, err := c.reminderRepo.SelectByUser(ctx, userID, req.IsDone, offset, req.PageSize)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取提醒列表失败")
		return
	}

	total, err := c.reminderRepo.CountByUser(ctx, userID, req.IsDone)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取提醒数量失败")
		return
	}

	response.SuccessResponse(ctx, "ok", gin.H{
		"list":  list,
		"total": total,
		"page":  req.Page,
	})
}

// TodayList 获取今日提醒
func (c *Controller) TodayList(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	if userID == 0 {
		userID = 1
	}

	list, err := c.reminderRepo.GetTodayReminders(ctx, userID)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取今日提醒失败")
		return
	}

	response.SuccessResponse(ctx, "ok", gin.H{
		"list": list,
	})
}

// PendingList 获取待处理提醒（已到期但未完成）
func (c *Controller) PendingList(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	if userID == 0 {
		userID = 1
	}

	list, err := c.reminderRepo.GetPendingReminders(ctx, userID)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取待处理提醒失败")
		return
	}

	response.SuccessResponse(ctx, "ok", gin.H{
		"list": list,
	})
}

// Stats 获取提醒统计
func (c *Controller) Stats(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	if userID == 0 {
		userID = 1
	}

	// 总数量
	total, _ := c.reminderRepo.CountByUser(ctx, userID, -1)
	// 待处理数量
	pending, _ := c.reminderRepo.CountByUser(ctx, userID, 0)
	// 今日提醒
	todayList, _ := c.reminderRepo.GetTodayReminders(ctx, userID)

	response.SuccessResponse(ctx, "ok", gin.H{
		"total":   total,
		"pending": pending,
		"today":   len(todayList),
	})
}

// MarkAsRead 标记为已读
func (c *Controller) MarkAsRead(ctx *gin.Context) {
	var req validates.ReminderIDValidate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateError(ctx, err, response.ValidateCommonError)
		return
	}

	if err := c.reminderRepo.MarkAsRead(ctx, req.ID); err != nil {
		response.ErrorResponse(ctx, response.DBUpdateCommonError, "标记已读失败")
		return
	}

	response.SuccessResponse(ctx, "标记成功", nil)
}

// MarkAsDone 标记为已完成
func (c *Controller) MarkAsDone(ctx *gin.Context) {
	var req validates.ReminderIDValidate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateError(ctx, err, response.ValidateCommonError)
		return
	}

	if err := c.reminderRepo.MarkAsDone(ctx, req.ID); err != nil {
		response.ErrorResponse(ctx, response.DBUpdateCommonError, "标记完成失败")
		return
	}

	response.SuccessResponse(ctx, "标记成功", nil)
}

// Delete 删除提醒
func (c *Controller) Delete(ctx *gin.Context) {
	// 支持从 query 参数获取 id
	idStr := ctx.Query("id")
	if idStr == "" {
		// 如果 query 没有，尝试从 JSON body 读取
		var req validates.ReminderIDValidate
		if err := ctx.ShouldBindJSON(&req); err != nil {
			response.ValidateError(ctx, err, response.ValidateCommonError)
			return
		}
		idStr = fmt.Sprintf("%d", req.ID)
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.ErrorResponse(ctx, response.ParamsCommonError, "无效的ID")
		return
	}

	if err := c.reminderRepo.Delete(ctx, id); err != nil {
		response.ErrorResponse(ctx, response.DBDeleteCommonError, "删除提醒失败")
		return
	}

	response.SuccessResponse(ctx, "删除成功", nil)
}
