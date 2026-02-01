package client

import (
	"omiai-server/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/iWuxc/go-wit/log"
)

func (c *Controller) Stats(ctx *gin.Context) {
	log.Infof("GET /api/client/stats requested")
	stats, err := c.Client.Stats(ctx)
	if err != nil {
		log.Errorf("Stats failed: %v", err)
		response.ErrorResponse(ctx, response.DBSelectCommonError, "获取统计信息失败")
		return
	}
	response.SuccessResponse(ctx, "ok", stats)
}
