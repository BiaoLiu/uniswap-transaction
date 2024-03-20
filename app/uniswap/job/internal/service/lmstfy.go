package service

import (
	"github.com/go-kratos/kratos/v2/log"

	"uniswap-transaction/app/uniswap/job/internal/biz"
)

type LmstfyService struct {
	log          *log.Helper
	ethereumRepo *biz.EthereumRepo
}

func NewLmstfyService(logger log.Logger, ethereumRepo *biz.EthereumRepo) *LmstfyService {
	return &LmstfyService{
		log:          log.NewHelper(logger),
		ethereumRepo: ethereumRepo,
	}
}
