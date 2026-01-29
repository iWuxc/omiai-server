//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"omiai-server/internal/controller"
	"omiai-server/internal/cron"
	"omiai-server/internal/data"
	"omiai-server/internal/data/omiai"
	"omiai-server/internal/middleware"
	"omiai-server/internal/queues"
	"omiai-server/internal/server"
	"omiai-server/internal/service"

	"github.com/google/wire"
	"github.com/iWuxc/go-wit/app"
)

func initApp(ctx context.Context) (*app.App, func(), error) {
	panic(wire.Build(
		cron.ProviderCronSet,
		queues.ProviderSet,
		middleware.ProviderMiddlewareSet,
		service.ProviderService,
		data.ProviderDataSet,
		server.ProviderServerSet,
		controller.ProviderController,
		omiai.ProviderOmiai,
		newApp))
}
