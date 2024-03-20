//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package biz

import (
	gconf "uniswap-transaction/protobuf/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	"uniswap-transaction/app/uniswap/job/internal/conf"
	"uniswap-transaction/app/uniswap/job/internal/data"
)

// initBizTest init bizTest.
func initBizTest(gconf.EnvMode, *conf.Data, *conf.Server, log.Logger) (*bizTest, func(), error) {
	panic(wire.Build(data.ProviderSet, ProviderSet, newBizTest))
}
