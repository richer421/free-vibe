//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"free-vibe-coding/internal/biz"
	"free-vibe-coding/internal/conf"
	"free-vibe-coding/internal/data"
	greeterdata "free-vibe-coding/internal/data/greeter"
	"free-vibe-coding/internal/server"
	"free-vibe-coding/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, greeterdata.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
