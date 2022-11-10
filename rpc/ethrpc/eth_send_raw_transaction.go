package ethrpc

import (
	"context"
	"encoding/json"

	"github.com/pandoprojects/pando-eth-rpc-adaptor/common"

	trpc "github.com/pandoprojects/pando/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_sendRawTransaction -----------------------------------

func (e *EthRPCService) SendRawTransaction(ctx context.Context, txBytes string) (result string, err error) {
	logger.Infof("eth_sendRawTransaction called")

	client := rpcc.NewRPCClient(common.GetPandoRPCEndpoint())
	rpcRes, rpcErr := client.Call("pando.BroadcastRawEthTransactionAsync", trpc.BroadcastRawTransactionAsyncArgs{TxBytes: txBytes})

	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.BroadcastRawTransactionAsyncResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		return trpcResult.TxHash, nil
	}

	resultIntf, err := common.HandlePandoRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		logger.Errorf("eth_sendRawTransaction, err: %v", err)
		return "", err
	}
	result = resultIntf.(string)

	logger.Infof("eth_sendRawTransaction, result: %v\n", result)

	return result, nil
}
