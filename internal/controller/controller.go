package controller

import (
	"omiai-server/internal/conf"
	"omiai-server/internal/controller/banner"
	"omiai-server/internal/controller/client"

	"github.com/google/wire"
)

var ProviderController = wire.NewSet(
	conf.GetConfig,
	banner.NewController,
	client.NewController,
)
