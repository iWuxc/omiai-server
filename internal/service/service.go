package service

import (
	"omiai-server/internal/service/banner"

	"github.com/google/wire"
)

var ProviderService = wire.NewSet(
	banner.NewService,
)
