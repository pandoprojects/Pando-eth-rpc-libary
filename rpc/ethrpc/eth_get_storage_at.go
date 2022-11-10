package ethrpc

import (
	"context"
	"encoding/json"
	"math"

	"github.com/pandoprojects/pando-eth-rpc-adaptor/common"

	trpc "github.com/pandoprojects/pando/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_getStorageAt -----------------------------------

func (e *EthRPCService) GetStorageAt(ctx context.Context, address string, storagePosition string, tag string) (result string, err error) {
	logger.Infof("eth_getStorageAt called")

	height := common.GetHeightByTag(tag)
	if height == math.MaxUint64 {
		height = 0 // 0 is interpreted as the last height by the pando.GetStorageAt method
	}

	client := rpcc.NewRPCClient(common.GetPandoRPCEndpoint())
	rpcRes, rpcErr := client.Call("pando.GetStorageAt", trpc.GetStorageAtArgs{
		Address:         address,
		StoragePosition: storagePosition,
		Height:          height})

	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.GetStorageAtResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		return trpcResult.Value, nil
	}

	resultIntf, err := common.HandlePandoRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		return "", err
	}

	result = resultIntf.(string)
	if result == "0000000000000000000000000000000000000000000000000000000000000000" {
		result = "0x0"
	}
	return result, nil
}
