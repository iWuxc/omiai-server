package dashboard

import (
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	client   biz_omiai.ClientInterface
	match    biz_omiai.MatchInterface
	reminder biz_omiai.ReminderInterface
}

func NewController(client biz_omiai.ClientInterface, match biz_omiai.MatchInterface, reminder biz_omiai.ReminderInterface) *Controller {
	return &Controller{
		client:   client,
		match:    match,
		reminder: reminder,
	}
}

type DashboardStats struct {
	ClientTotal     int64 `json:"client_total"`
	ClientToday     int64 `json:"client_today"`
	ClientMonth     int64 `json:"client_month"`
	MatchTotal      int64 `json:"match_total"`
	MatchMonth      int64 `json:"match_month"`
	FollowUpPending int64 `json:"follow_up_pending"`
}

type TodoItem struct {
	ID         int64  `json:"id"`
	Type       string `json:"type"`
	Title      string `json:"title"`
	Priority   string `json:"priority"`
	Status     string `json:"status"`
	ClientID   *int64 `json:"client_id,omitempty"`
	ClientName string `json:"client_name,omitempty"`
	CreatedAt  string `json:"created_at"`
}

func (c *Controller) Stats(ctx *gin.Context) {
	stats, err := c.client.GetDashboardStats(ctx)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取统计数据失败")
		return
	}

	// 获取匹配统计
	matchStats, err := c.match.Stats(ctx)
	if err == nil && matchStats != nil {
		if totalMatches, ok := matchStats["total_matches"].(int64); ok {
			stats["match_total"] = totalMatches
		}

		// 本月匹配
		// TODO: 实现本月匹配统计逻辑
	}

	// 获取待办提醒
	// TODO: Get user ID from context. Assuming default 1 for now if not set.
	userID := ctx.GetUint64("user_id")
	if userID == 0 {
		userID = 1
	}
	
	pendingReminders, err := c.reminder.GetPendingReminders(userID)
	if err == nil {
		stats["follow_up_pending"] = int64(len(pendingReminders))
	}

	response.SuccessResponse(ctx, "ok", stats)
}

func (c *Controller) GetTodos(ctx *gin.Context) {
	// 获取当前用户ID
	userID := ctx.GetUint64("user_id")
	if userID == 0 {
		userID = 1 // 默认用户
	}

	// 从提醒系统获取待办事项
	pendingTasks, err := c.reminder.GetPendingReminders(userID)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取待办事项失败")
		return
	}

	todos := make([]TodoItem, 0)
	for _, task := range pendingTasks {
		// Priority logic could be complex based on Rule or Content keywords
		priority := "medium" 
		taskType := "reminder"
		
		// Simple keyword matching for type and priority demo
		// In production, ReminderTask should probably have Type and Priority fields
		
		todos = append(todos, TodoItem{
			ID:        task.ID,
			Type:      taskType,
			Title:     task.Content,
			Priority:  priority,
			Status:    task.Status,
			ClientID:  &task.ClientID,
			// ClientName: fetch client name if needed, or join in repo
			CreatedAt: task.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	response.SuccessResponse(ctx, "ok", todos)
}
