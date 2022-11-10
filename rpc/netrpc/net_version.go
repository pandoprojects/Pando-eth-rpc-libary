package netrpc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pandotoken/pando-eth-rpc-adaptor/common"
	hexutil "github.com/pandotoken/pando/common/hexutil"
	"github.com/pandotoken/pando/ledger/types"
	trpc "github.com/pandotoken/pando/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

type chainIDResultWrapper struct {
	chainID string
}

// ------------------------------- net_version -----------------------------------

func (e *NetRPCService) Version(ctx context.Context) (result string, err error) {
	logger.Infof("net_version called")

	client := rpcc.NewRPCClient(common.GetPandoRPCEndpoint())
	rpcRes, rpcErr := client.Call("pando.GetStatus", trpc.GetStatusArgs{})
	var blockHeight uint64
	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.GetStatusResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		re := chainIDResultWrapper{
			chainID: trpcResult.ChainID,
		}
		blockHeight = uint64(trpcResult.LatestFinalizedBlockHeight)
		return re, nil
	}

	resultIntf, err := common.HandlePandoRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		return "", err
	}
	pandoChainIDResult, ok := resultIntf.(chainIDResultWrapper)
	if !ok {
		return "", fmt.Errorf("failed to convert chainIDResultWrapper")
	}

	pandoChainID := pandoChainIDResult.chainID
	ethChainID := types.MapChainID(pandoChainID, blockHeight).Uint64()
	result = hexutil.EncodeUint64(ethChainID)

	return result, nil
}
