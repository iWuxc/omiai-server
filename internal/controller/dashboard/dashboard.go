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

// GetInteractionLeads 获取C端高意向线索
func (c *Controller) GetInteractionLeads(ctx *gin.Context) {
	managerID := ctx.GetUint64("user_id")
	if managerID == 0 {
		response.ErrorResponse(ctx, response.AuthCommonError, "未授权")
		return
	}

	// 限制返回最新 20 条线索
	leads, err := c.client.GetInteractionLeads(ctx, managerID, 0, 20)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取线索失败")
		return
	}

	// 组装返回数据，带上客户详情
	var result []map[string]interface{}
	for _, lead := range leads {
		fromClient, _ := c.client.Get(ctx, lead.FromClientID)
		toClient, _ := c.client.Get(ctx, lead.ToClientID)

		if fromClient != nil && toClient != nil {
			result = append(result, map[string]interface{}{
				"interaction_id": lead.ID,
				"from_client": map[string]interface{}{
					"id":     fromClient.ID,
					"name":   fromClient.Name,
					"avatar": fromClient.Avatar,
					"phone":  fromClient.Phone,
				},
				"to_client": map[string]interface{}{
					"id":     toClient.ID,
					"name":   toClient.Name,
					"avatar": toClient.Avatar,
					"phone":  toClient.Phone,
				},
				"created_at": lead.CreatedAt.Format("2006-01-02 15:04:05"),
			})
		}
	}

	response.SuccessResponse(ctx, "ok", result)
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
	// 判断角色是否为店长(admin)，如果是，可以看整个机构的，否则只能看自己的
	role := ctx.GetString("role")
	managerID := ctx.GetUint64("user_id")

	// 通过 Context 将筛选条件传给 Repo 层
	// 如果是普通红娘，可以在 Repo 里面增加 manager_id 过滤
	ctx.Set("filter_role", role)
	ctx.Set("filter_manager_id", managerID)

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
			ID:       task.ID,
			Type:     taskType,
			Title:    task.Content,
			Priority: priority,
			Status:   task.Status,
			ClientID: &task.ClientID,
			// ClientName: fetch client name if needed, or join in repo
			CreatedAt: task.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	response.SuccessResponse(ctx, "ok", todos)
}

// UpdateInteractionStatus 更新线索状态
func (c *Controller) UpdateInteractionStatus(ctx *gin.Context) {
	var req struct {
		InteractionID uint64 `json:"interaction_id" binding:"required"`
		Status        int8   `json:"status" binding:"required,oneof=1 2"` // 1已跟进 2已忽略
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateError(ctx, err, response.ValidateCommonError)
		return
	}

	// 这里需要单独查询 interaction 并更新状态
	if err := c.client.UpdateInteractionStatus(ctx, req.InteractionID, req.Status); err != nil {
		response.ErrorResponse(ctx, response.DBUpdateCommonError, "状态更新失败")
		return
	}

	response.SuccessResponse(ctx, "状态更新成功", nil)
}
