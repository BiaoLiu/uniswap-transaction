package server

import (
	"context"

	"uniswap-transaction/app/uniswap/job/internal/conf"
	"uniswap-transaction/app/uniswap/job/internal/service"
	"uniswap-transaction/pkg/alert"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type EthereumServer struct {
	conf         *conf.Data
	alertManager *alert.Manager
	client       *ethclient.Client
	doneChan     chan struct{}

	ethereumService *service.EthereumService
}

func NewEthereumServer(c *conf.Data, logger log.Logger, alertManager *alert.Manager, ethereumService *service.EthereumService) *EthereumServer {
	client, err := ethclient.Dial(c.Ethereum.ApiUrl)
	if err != nil {
		log.Fatal(err)
	}
	return &EthereumServer{
		conf:            c,
		alertManager:    alertManager,
		client:          client,
		doneChan:        make(chan struct{}),
		ethereumService: ethereumService,
	}
}

func (s *EthereumServer) Start(ctx context.Context) error {
	log.Infof("[Job] ethereum server start")
	err := s.SubscribeNewBlock(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *EthereumServer) Stop(ctx context.Context) error {
	log.Infof("[Job] ethereum server stopping")
	close(s.doneChan)
	return nil
}

func (s *EthereumServer) SubscribeNewBlock(ctx context.Context) error {
	log.Infof("[Job] subscribe new blocks start")

	headers := make(chan *types.Header)
	sub, err := s.client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		return errors.Wrapf(err, "subscribe new block error")
	}

	eg, _ := errgroup.WithContext(ctx)
	eg.Go(func() error {
		for {
			select {
			case <-s.doneChan:
				close(headers)
				return errors.New("stop subscribe new block")
			case err := <-sub.Err():
				close(headers)
				return err
			case header := <-headers:
				// fmt.Printf("%+v\n", header)
				err = s.ethereumService.HandleTransaction(ctx, header.Number)
				if err != nil {
					log.Errorf("handle transaction error. hash:%v", header.Hash().Hex())
					continue
				}
			}
		}
	},
	)
	return nil
}
