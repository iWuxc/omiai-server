package server

import (
	"omiai-server/internal/api"

	"github.com/google/wire"
)

// ProviderServerSet is server providers.
var ProviderServerSet = wire.NewSet(
	NewGinEngine, NewHTTPServer, // server
	wire.Struct(new(Router), "*"),
	wire.Bind(new(api.RouterInterface), new(*Router)),
	NewRedisX,
	NewRedisSync,
)
