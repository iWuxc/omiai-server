package service

import (
	"omiai-server/internal/service/banner"
	"omiai-server/internal/service/chat_parser"

	"github.com/google/wire"
)

var ProviderService = wire.NewSet(
	banner.NewService,
	chat_parser.NewChatParser,
)
