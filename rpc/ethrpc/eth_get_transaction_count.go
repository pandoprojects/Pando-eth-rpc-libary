package ethrpc

import (
	"context"
	"encoding/json"
	"math"

	"github.com/pandoprojects/pando-eth-rpc-adaptor/common"
	hexutil "github.com/pandoprojects/pando/common/hexutil"
	"github.com/pandoprojects/pando/ledger/types"
	trpc "github.com/pandoprojects/pando/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_getTransactionCount -----------------------------------

func (e *EthRPCService) GetTransactionCount(ctx context.Context, address string, tag string) (result string, err error) {
	logger.Infof("eth_getTransactionCount called, address: %v, tag: %v", address, tag)
	height := common.GetHeightByTag(tag)
	if height == math.MaxUint64 {
		height = 0 // 0 is interpreted as the last height by the pando.GetAccount method
	}

	client := rpcc.NewRPCClient(common.GetPandoRPCEndpoint())
	rpcRes, rpcErr := client.Call("pando.GetAccount", trpc.GetAccountArgs{Address: address, Height: height, Preview: true})

	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.GetAccountResult{Account: &types.Account{}}
		json.Unmarshal(jsonBytes, &trpcResult)
		return trpcResult.Account.Sequence, nil
	}

	resultIntf, err := common.HandlePandoRPCResponse(rpcRes, rpcErr, parse)

	if err != nil {
		return "0x0", nil
	}

	// result = fmt.Sprintf("0x%x", resultIntf.(*big.Int))
	result = hexutil.EncodeUint64(resultIntf.(uint64))

	return result, nil
}
