package biz

import (
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewEthereumRepo,
)

type bizTest struct {
	ethereumRepo *EthereumRepo
}

func newBizTest(ethereumRepo *EthereumRepo) *bizTest {
	return &bizTest{
		ethereumRepo: ethereumRepo,
	}
}
