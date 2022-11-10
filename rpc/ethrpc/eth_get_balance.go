package ethrpc

import (
	"context"
	"encoding/json"
	"math"
	"math/big"

	"github.com/pandoprojects/pando-eth-rpc-adaptor/common"
	"github.com/pandoprojects/pando/ledger/types"
	trpc "github.com/pandoprojects/pando/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_getBalance -----------------------------------

func (e *EthRPCService) GetBalance(ctx context.Context, address string, tag string) (result string, err error) {
	logger.Infof("eth_getBalance called")

	height := common.GetHeightByTag(tag)
	if height == math.MaxUint64 {
		height = 0 // 0 is interpreted as the last height by the pando.GetAccount method
	}

	client := rpcc.NewRPCClient(common.GetPandoRPCEndpoint())
	rpcRes, rpcErr := client.Call("pando.GetAccount", trpc.GetAccountArgs{Address: address, Height: height})

	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.GetAccountResult{Account: &types.Account{}}
		json.Unmarshal(jsonBytes, &trpcResult)
		return trpcResult.Account.Balance.PTXWei, nil
	}

	resultIntf, err := common.HandlePandoRPCResponse(rpcRes, rpcErr, parse)

	if err != nil {
		return "0x0", nil
	}

	// result = fmt.Sprintf("0x%x", resultIntf.(*big.Int))
	result = "0x" + (resultIntf.(*big.Int)).Text(16)

	return result, nil
}
