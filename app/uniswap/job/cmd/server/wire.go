//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	"uniswap-transaction/app/uniswap/job/internal/biz"
	"uniswap-transaction/app/uniswap/job/internal/conf"
	"uniswap-transaction/app/uniswap/job/internal/data"
	"uniswap-transaction/app/uniswap/job/internal/server"
	"uniswap-transaction/app/uniswap/job/internal/service"
	gconf "uniswap-transaction/protobuf/conf"
)

// initApp init kratos application.
func initApp(gconf.EnvMode, *conf.Server, *conf.Data, *conf.Alert, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
