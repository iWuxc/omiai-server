package omiai

import (
	"github.com/google/wire"
	"omiai-server/internal/data/event"
	"omiai-server/internal/data/tenant"
)

var ProviderData = wire.NewSet(
	NewBannerRepo,
	NewChinaRegionRepo,
	NewClientRepo,
	NewMatchRepo,
	NewReminderRepo,
	NewTemplateRepo,
	NewUserRepo,
	NewAIMatchRepo,
	tenant.NewTenantRepo,
	event.NewEventRepo,
)
