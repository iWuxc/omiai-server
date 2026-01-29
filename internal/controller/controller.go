package controller

import (
	"omiai-server/internal/conf"
	"omiai-server/internal/controller/banner"

	"github.com/google/wire"
)

var ProviderController = wire.NewSet(
	conf.GetConfig,
	banner.NewController,
)
