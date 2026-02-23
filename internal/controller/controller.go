package controller

import (
	"omiai-server/internal/conf"
	"omiai-server/internal/controller/ai"
	"omiai-server/internal/controller/auth"
	"omiai-server/internal/controller/banner"
	"omiai-server/internal/controller/china_region"
	"omiai-server/internal/controller/client"
	"omiai-server/internal/controller/common"
	"omiai-server/internal/controller/dashboard"
	"omiai-server/internal/controller/match"
	"omiai-server/internal/controller/reminder"
	"omiai-server/internal/controller/template"

	"github.com/google/wire"
)

var ProviderController = wire.NewSet(
	conf.GetConfig,
	ai.NewController,
	auth.NewController,
	banner.NewController,
	china_region.NewController,
	client.NewController,
	common.NewController,
	dashboard.NewController,
	match.NewController,
	reminder.NewController,
	template.NewController,
)
