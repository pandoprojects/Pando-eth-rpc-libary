package ethrpc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pandoprojects/pando-eth-rpc-adaptor/common"

	"github.com/pandoprojects/pando/common/hexutil"
	trpc "github.com/pandoprojects/pando/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

type syncingResultWrapper struct {
	*common.EthSyncingResult
	syncing bool
}

// ------------------------------- eth_syncing -----------------------------------
func (e *EthRPCService) Syncing(ctx context.Context) (result interface{}, err error) {
	logger.Infof("eth_syncing called")
	client := rpcc.NewRPCClient(common.GetPandoRPCEndpoint())
	rpcRes, rpcErr := client.Call("pando.GetStatus", trpc.GetStatusArgs{})
	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.GetStatusResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		re := syncingResultWrapper{&common.EthSyncingResult{}, false}
		re.syncing = trpcResult.Syncing
		if trpcResult.Syncing {
			re.StartingBlock = 1
			re.CurrentBlock = hexutil.Uint64(trpcResult.CurrentHeight)
			re.HighestBlock = hexutil.Uint64(trpcResult.LatestFinalizedBlockHeight)
			re.PulledStates = re.CurrentBlock
			re.KnownStates = re.CurrentBlock
		}
		return re, nil
	}

	resultIntf, err := common.HandlePandoRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		return "", err
	}
	pandoSyncingResult, ok := resultIntf.(syncingResultWrapper)
	if !ok {
		return nil, fmt.Errorf("failed to convert syncingResultWrapper")
	}
	if !pandoSyncingResult.syncing {
		result = false
	} else {
		result = pandoSyncingResult.EthSyncingResult
	}

	return result, nil
}
