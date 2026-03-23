package omiai

import (
	"github.com/google/wire"
	"omiai-server/internal/data/tenant"
)

var ProviderOmiai = wire.NewSet(
	NewBannerRepo,
	NewClientRepo,
	NewMatchRepo,
	NewUserRepo,
	NewReminderRepo,
	NewChinaRegionRepo,
	NewTemplateRepo,
	NewAIMatchRepo,
	tenant.NewTenantRepo,
)
