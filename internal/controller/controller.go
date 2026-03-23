package controller

import (
	"omiai-server/internal/conf"
	"omiai-server/internal/controller/ai"
	"omiai-server/internal/controller/auth"
	"omiai-server/internal/controller/banner"
	"omiai-server/internal/controller/c_auth"
	"omiai-server/internal/controller/c_client"
	"omiai-server/internal/controller/c_interact"
	"omiai-server/internal/controller/c_recommend"
	"omiai-server/internal/controller/china_region"
	"omiai-server/internal/controller/client"
	"omiai-server/internal/controller/common"
	"omiai-server/internal/controller/dashboard"
	"omiai-server/internal/controller/match"
	"omiai-server/internal/controller/reminder"
	"omiai-server/internal/controller/template"
	"omiai-server/internal/service/notification"

	"github.com/google/wire"
)

var ProviderController = wire.NewSet(
	conf.GetConfig,
	notification.NewNotificationService,
	ai.NewController,
	auth.NewController,
	c_auth.NewController,
	c_client.NewController,
	c_recommend.NewController,
	c_interact.NewController,
	banner.NewController,
	china_region.NewController,
	client.NewController,
	common.NewController,
	dashboard.NewController,
	match.NewController,
	reminder.NewController,
	template.NewController,
)
