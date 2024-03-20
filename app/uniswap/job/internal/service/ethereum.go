package service

import (
	"context"
	"math/big"

	"uniswap-transaction/app/uniswap/job/internal/data"

	"uniswap-transaction/app/uniswap/job/internal/biz"
	"uniswap-transaction/app/uniswap/job/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

type CreateTransactionReq struct {
	Hash       string
	Token0     string
	Token1     string
	Type       string
	Amount0    int64
	Amount1    int64
	Amount0In  int64
	Amount1Out int64
	From       string
	To         string
}

type EthereumService struct {
	log          *log.Helper
	conf         *conf.Data
	ethereumRepo *biz.EthereumRepo
}

func NewEthereumService(logger log.Logger, conf *conf.Data, ethereumRepo *biz.EthereumRepo) *EthereumService {
	return &EthereumService{
		log:          log.NewHelper(logger),
		conf:         conf,
		ethereumRepo: ethereumRepo,
	}
}

// HandleTransaction handle transaction
func (s *EthereumService) HandleTransaction(ctx context.Context, blockNumber *big.Int) error {
	err := s.ethereumRepo.HandleTransaction(ctx, blockNumber)
	return err
}

// CreateTransaction create transaction
func (s *EthereumService) CreateTransaction(ctx context.Context, req *CreateTransactionReq) (exist bool, err error) {
	transaction := &data.Transaction{
		Hash:       req.Hash,
		Token0:     req.Token0,
		Token1:     req.Token1,
		Amount0:    req.Amount0,
		Amount1:    req.Amount1,
		Amount0In:  req.Amount0In,
		Amount1Out: req.Amount1Out,
		From:       req.From,
		To:         req.To,
	}
	exist, err = s.ethereumRepo.CreateTransaction(ctx, transaction)
	return exist, err
}
