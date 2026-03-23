package event

import (
	biz_event "omiai-server/internal/biz/event"
	"omiai-server/pkg/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	Event biz_event.EventInterface
}

func NewController(event biz_event.EventInterface) *Controller {
	return &Controller{Event: event}
}

type PublishEventRequest struct {
	Title       string `json:"title" binding:"required"`
	Cover       string `json:"cover"`
	Description string `json:"description"`
	Address     string `json:"address" binding:"required"`
	StartTime   string `json:"start_time" binding:"required"`
	EndTime     string `json:"end_time" binding:"required"`
	PriceCoins  int    `json:"price_coins" binding:"min=0"`
	MaxQuota    int    `json:"max_quota" binding:"required,min=2"`
}

// Publish 发布线下活动 (B端红娘/店长)
func (c *Controller) Publish(ctx *gin.Context) {
	var req PublishEventRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateError(ctx, err, response.ValidateCommonError)
		return
	}

	layout := "2006-01-02 15:04:05"
	startTime, _ := time.Parse(layout, req.StartTime)
	endTime, _ := time.Parse(layout, req.EndTime)

	event := &biz_event.Event{
		Title:       req.Title,
		Cover:       req.Cover,
		Description: req.Description,
		Address:     req.Address,
		StartTime:   startTime,
		EndTime:     endTime,
		PriceCoins:  req.PriceCoins,
		MaxQuota:    req.MaxQuota,
		Status:      1, // 报名中
	}

	if err := c.Event.Create(ctx, event); err != nil {
		response.ErrorResponse(ctx, response.DBInsertCommonError, "发布活动失败")
		return
	}

	response.SuccessResponse(ctx, "发布成功", event)
}

// List 获取活动列表 (B端管理)
func (c *Controller) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	status, _ := strconv.Atoi(ctx.DefaultQuery("status", "0"))
	
	offset := (page - 1) * size

	list, total, err := c.Event.List(ctx, int8(status), offset, size)
	if err != nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取列表失败")
		return
	}

	response.SuccessResponse(ctx, "success", map[string]interface{}{
		"list":  list,
		"total": total,
	})
}
