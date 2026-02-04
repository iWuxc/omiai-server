package client

import (
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/data"
	"omiai-server/internal/service/chat_parser"
)

type Controller struct {
	db                *data.DB
	client            biz_omiai.ClientInterface
	chatParserService *chat_parser.ChatParser
}

func NewController(db *data.DB, client biz_omiai.ClientInterface, chatParserService *chat_parser.ChatParser) *Controller {
	return &Controller{db: db, client: client, chatParserService: chatParserService}
}
