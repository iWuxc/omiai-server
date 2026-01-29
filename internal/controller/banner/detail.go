package banner

import (
	"omiai-server/pkg/response"

	"github.com/gin-gonic/gin"
)

func (c *Controller) Detail(ctx *gin.Context) {
	detail := c.bannerService.GetDisplayAttributes()
	response.SuccessResponse(ctx, "ok", detail)
}
