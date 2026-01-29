//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"omiai-server/cmd/script/command"
	"omiai-server/internal/data"

	"github.com/google/wire"
)

func initApp(ctx context.Context) (*InitCmd, func(), error) {
	panic(wire.Build(
		ProviderSet,
		command.ProviderSet,
		data.ProviderDataSet,
		//service.ProviderService,
		//server.ProviderServerSet,
		//omiai.ProviderOmiai,
	))
}
