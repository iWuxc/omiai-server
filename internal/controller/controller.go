package controller

import (
	"omiai-server/internal/conf"
	"omiai-server/internal/controller/auth"
	"omiai-server/internal/controller/banner"
	"omiai-server/internal/controller/client"
	"omiai-server/internal/controller/common"
	"omiai-server/internal/controller/match"

	"github.com/google/wire"
)

var ProviderController = wire.NewSet(
	conf.GetConfig,
	auth.NewController,
	banner.NewController,
	client.NewController,
	common.NewController,
	match.NewController,
)
