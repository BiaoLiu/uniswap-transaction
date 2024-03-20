package biz

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"

	"uniswap-transaction/app/uniswap/job/internal/pkg/contract"
)

func TestUnpack(t *testing.T) {
	const hexdata = `0x7ff36ab5000000000000000000000000000000000000000000000000001a6dbd51f512690000000000000000000000000000000000000000000000000000000000000080000000000000000000000000d8587c02d9891a4ef71ea7e6004e2ce6ddffd10a0000000000000000000000000000000000000000000000000000000061dbd5a70000000000000000000000000000000000000000000000000000000000000002000000000000000000000000bb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c000000000000000000000000c146b7cdbaff065090077151d391f4c96aa09e0c`

	abi, err := abi.JSON(strings.NewReader(routerABIJSON))
	assert.NoError(t, err)

	data, err := hex.DecodeString(hexdata[10:])
	assert.NoError(t, err)

	params := make(map[string]interface{})
	err = abi.Methods["swapExactETHForTokens"].Inputs.UnpackIntoMap(params, data)
	assert.NoError(t, err)

	fmt.Printf("%+v", params)
	amountOutMin := params["amountOutMin"].(*big.Int)
	path := params["path"].([]common.Address)
	to := params["to"].(common.Address)
	deadline := params["deadline"].(*big.Int)
	fmt.Println(amountOutMin, path, to, deadline)
}

func TestQueryTransaction(t *testing.T) {
	blockHash := common.HexToHash("0x90f8b94a44822270ca5805e0f6c3291e952ef4001592bbd568bb08e36b2a236d")
	block, err := biz.ethereumRepo.data.EthereumClient.BlockByHash(context.Background(), blockHash)
	assert.NoError(t, err)

	routerABI, err := abi.JSON(strings.NewReader(routerABIJSON))
	assert.NoError(t, err)

	pairABI, err := abi.JSON(strings.NewReader(pairABIJSON))
	assert.NoError(t, err)

	for _, tx := range block.Transactions() {
		if tx.To() == nil || tx.To().Hex() != biz.ethereumRepo.conf.Ethereum.UniswapRouterAddress {
			continue
		}

		input := strings.TrimPrefix(string(tx.Data()), "0x")
		inputBytes := []byte(input)

		methodID := inputBytes[:4]
		inputBytes = inputBytes[4:]
		method, err := routerABI.MethodById(methodID)
		if err != nil {
			continue
		}

		methodName := method.Name
		if !strings.HasPrefix(methodName, "addLiquidity") &&
			!strings.HasPrefix(methodName, "removeLiquidity") &&
			!strings.HasPrefix(methodName, "swap") {
			continue
		}

		receipt, err := biz.ethereumRepo.data.EthereumClient.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			fmt.Println(err)
			continue
		}
		swapSignHash := crypto.Keccak256Hash(swapSignature)

		var senderEvent SwapEvent
		for _, vLog := range receipt.Logs {
			switch vLog.Topics[0].Hex() {
			case swapSignHash.Hex():
				err := pairABI.UnpackIntoInterface(&senderEvent, "Swap", vLog.Data)
				if err != nil {
					fmt.Println(err)
					continue
				}
				uniswap, err := contract.NewUniswap(vLog.Address, biz.ethereumRepo.data.EthereumClient)
				if err != nil {
					fmt.Println(err)
				}
				test1, err := uniswap.Token0(nil)
				test2, err := uniswap.Name(nil)
				fmt.Println(test1)
				fmt.Println(test2)
			}
			fmt.Println(vLog)
		}

		params := make(map[string]interface{})
		err = routerABI.Methods[methodName].Inputs.UnpackIntoMap(params, inputBytes)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("%+v", params)
		amountOutMin := params["amountOutMin"].(*big.Int)
		path := params["path"].([]common.Address)
		to := params["to"].(common.Address)
		deadline := params["deadline"].(*big.Int)

		fmt.Println(amountOutMin, path, to, deadline)
	}
}

func TestHandleTransaction(t *testing.T) {
	// blockHash := common.HexToHash("0x90f8b94a44822270ca5805e0f6c3291e952ef4001592bbd568bb08e36b2a236d")
	// err := biz.ethereumRepo.HandleTransaction(context.Background(), blockHash)
	number := big.NewInt(19475148)
	err := biz.ethereumRepo.HandleTransaction(context.Background(), number)
	assert.NoError(t, err)
}
