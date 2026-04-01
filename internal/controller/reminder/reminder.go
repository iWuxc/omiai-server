package reminder

import (
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Rule Handlers
func (c *Controller) CreateRule(ctx *gin.Context) {
	var req biz_omiai.AutoReminderRule
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(ctx, response.ValidateCommonError, err.Error())
		return
	}
	if err := c.reminderRepo.CreateRule(&req); err != nil {
		response.ErrorResponse(ctx, response.DBInsertCommonError, "创建规则失败")
		return
	}
	response.SuccessResponse(ctx, "创建成功", req)
}

func (c *Controller) ListRules(ctx *gin.Context) {
	rules, err := c.reminderRepo.ListRules()
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取规则失败")
		return
	}
	response.SuccessResponse(ctx, "获取成功", rules)
}

// Task Handlers
func (c *Controller) List(ctx *gin.Context) {
	var req struct {
		Status string `form:"status"`
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ErrorResponse(ctx, response.ValidateCommonError, err.Error())
		return
	}

	var tasks []*biz_omiai.ReminderTask
	var err error

	if req.Status != "" {
		tasks, err = c.reminderRepo.GetTasksByStatus(req.Status)
	} else {
		tasks, err = c.reminderRepo.GetAllTasks()
	}

	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取提醒列表失败")
		return
	}
	response.SuccessResponse(ctx, "获取成功", tasks)
}

func (c *Controller) TodayList(ctx *gin.Context) {
	userID := uint64(0)
	tasks, err := c.reminderRepo.GetTodayReminders(userID)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取今日提醒失败")
		return
	}
	response.SuccessResponse(ctx, "获取成功", tasks)
}

func (c *Controller) ListPendingTasks(ctx *gin.Context) {
	userID := uint64(0)
	tasks, err := c.reminderRepo.GetPendingReminders(userID)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取待办任务失败")
		return
	}
	response.SuccessResponse(ctx, "获取成功", tasks)
}

func (c *Controller) MarkAsRead(ctx *gin.Context) {
	var req struct {
		ID int64 `json:"id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(ctx, response.ValidateCommonError, err.Error())
		return
	}
	if err := c.reminderRepo.MarkAsRead(req.ID); err != nil {
		response.ErrorResponse(ctx, response.DBUpdateCommonError, "标记已读失败")
		return
	}
	response.SuccessResponse(ctx, "操作成功", nil)
}

func (c *Controller) MarkAsDone(ctx *gin.Context) {
	var req struct {
		ID int64 `json:"id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(ctx, response.ValidateCommonError, err.Error())
		return
	}
	if err := c.reminderRepo.MarkAsDone(req.ID); err != nil {
		response.ErrorResponse(ctx, response.DBUpdateCommonError, "标记完成失败")
		return
	}
	response.SuccessResponse(ctx, "操作成功", nil)
}

func (c *Controller) CompleteTask(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err := c.reminderRepo.CompleteTask(id); err != nil {
		response.ErrorResponse(ctx, response.DBUpdateCommonError, "操作失败")
		return
	}
	response.SuccessResponse(ctx, "操作成功", nil)
}

func (c *Controller) Delete(ctx *gin.Context) {
	var req struct {
		ID int64 `form:"id" binding:"required"`
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ErrorResponse(ctx, response.ValidateCommonError, err.Error())
		return
	}
	if err := c.reminderRepo.Delete(req.ID); err != nil {
		response.ErrorResponse(ctx, response.DBDeleteCommonError, "删除失败")
		return
	}
	response.SuccessResponse(ctx, "删除成功", nil)
}

// CheckAndGenerateTasks 定时任务入口：检查规则并生成任务
// 注意：这个接口应该由定时任务系统（如 cron）调用，或者做成后台常驻协程
func (c *Controller) CheckAndGenerateTasks(ctx *gin.Context) {
	// 简易实现：遍历所有规则，查找符合条件的 Client，生成 Task
	// 实际生产环境应使用更高效的查询或事件驱动
	rules, _ := c.reminderRepo.ListRules()
	count := 0

	for _, rule := range rules {
		if !rule.IsEnabled {
			continue
		}
		// TODO: 这里需要根据 rule.TriggerType 和 TriggerCondition 查询 Client
		// 暂时仅演示框架，实际逻辑需要复杂的 SQL 构建器
		// 示例：TriggerType="NoContact", DelayDays=7 -> 查找 last_contact_at < now - 7 days
	}

	response.SuccessResponse(ctx, "任务生成完成", gin.H{"count": count})
}

// Stats 获取提醒统计数据
func (c *Controller) Stats(ctx *gin.Context) {
	// 获取当前用户ID (从上下文获取，这里简化处理为0表示全部)
	userID := uint64(0)

	pendingCount, err := c.reminderRepo.CountByUser(userID, 0)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取待办统计失败")
		return
	}

	completedCount, err := c.reminderRepo.CountByUser(userID, 1)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取已完成统计失败")
		return
	}

	allCount, err := c.reminderRepo.CountByUser(userID, -1)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取总数统计失败")
		return
	}

	response.SuccessResponse(ctx, "获取成功", gin.H{
		"pending":   pendingCount,
		"completed": completedCount,
		"total":     allCount,
	})
}
