package ethrpc

import (
	"context"
	"encoding/json"

	"github.com/pandoprojects/pando-eth-rpc-adaptor/common"

	trpc "github.com/pandoprojects/pando/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_protocolVersion -----------------------------------

func (e *EthRPCService) ProtocolVersion(ctx context.Context) (result string, err error) {
	logger.Infof("eth_protocolVersion called")

	client := rpcc.NewRPCClient(common.GetPandoRPCEndpoint())
	rpcRes, rpcErr := client.Call("pando.GetVersion", trpc.GetVersionArgs{})

	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.GetVersionResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		return trpcResult.Version, nil
	}

	resultIntf, err := common.HandlePandoRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		return "", err
	}
	result = resultIntf.(string)

	return result, nil
}
