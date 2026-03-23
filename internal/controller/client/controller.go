package client

import (
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/data"
	"omiai-server/internal/service/chat_parser"
	"omiai-server/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iWuxc/go-wit/log"
)

type Controller struct {
	db                *data.DB
	client            biz_omiai.ClientInterface
	chatParserService *chat_parser.ChatParser
}

func NewController(db *data.DB, client biz_omiai.ClientInterface, chatParserService *chat_parser.ChatParser) *Controller {
	return &Controller{db: db, client: client, chatParserService: chatParserService}
}

// Verify C端资料审核通过
func (c *Controller) Verify(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.ErrorResponse(ctx, response.ParamsCommonError, "参数错误")
		return
	}

	client, err := c.client.Get(ctx, id)
	if err != nil || client == nil {
		response.ErrorResponse(ctx, response.DBSelectCommonError, "客户不存在")
		return
	}

	client.IsVerified = true
	if err := c.client.Update(ctx, client); err != nil {
		log.Errorf("Verify client failed: %v", err)
		response.ErrorResponse(ctx, response.DBUpdateCommonError, "审核失败")
		return
	}

	response.SuccessResponse(ctx, "审核成功", nil)
}
