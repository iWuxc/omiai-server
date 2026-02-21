package dashboard

import (
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	client biz_omiai.ClientInterface
	match  biz_omiai.MatchInterface
}

func NewController(client biz_omiai.ClientInterface, match biz_omiai.MatchInterface) *Controller {
	return &Controller{
		client: client,
		match:  match,
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

	response.SuccessResponse(ctx, "ok", stats)
}

func (c *Controller) GetTodos(ctx *gin.Context) {
	// TODO: 从提醒系统获取待办事项
	// 暂时返回示例数据
	todos := []TodoItem{
		{
			ID:       1,
			Type:     "follow_up",
			Title:    "客户张三需要跟进",
			Priority: "high",
			Status:   "pending",
		},
		{
			ID:       2,
			Type:     "birthday",
			Title:    "客户李四今天生日",
			Priority: "medium",
			Status:   "pending",
		},
	}

	response.SuccessResponse(ctx, "ok", todos)
}
