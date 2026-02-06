package omiai

import "github.com/google/wire"

var ProviderOmiai = wire.NewSet(
	NewBannerRepo,
	NewClientRepo,
	NewMatchRepo,
	NewUserRepo,
	NewReminderRepo,
)
