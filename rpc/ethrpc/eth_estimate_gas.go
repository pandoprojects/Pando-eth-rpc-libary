package ethrpc

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/spf13/viper"
	"github.com/pandoprojects/pando-eth-rpc-adaptor/common"
	tcommon "github.com/pandoprojects/pando/common"
	hexutil "github.com/pandoprojects/pando/common/hexutil"

	trpc "github.com/pandoprojects/pando/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_estimateGas -----------------------------------

func (e *EthRPCService) EstimateGas(ctx context.Context, argObj common.EthSmartContractArgObj) (result string, err error) {
	logger.Infof("eth_estimateGas called")

	sctxBytes, err := common.GetSctxBytes(argObj)
	if err != nil {
		logger.Errorf("eth_estimateGas: Failed to get smart contract bytes: %+v\n", argObj)
		return result, err
	}

	client := rpcc.NewRPCClient(common.GetPandoRPCEndpoint())

	rpcRes, rpcErr := client.Call("pando.CallSmartContract", trpc.CallSmartContractArgs{SctxBytes: hex.EncodeToString(sctxBytes)})

	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.CallSmartContractResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		if len(trpcResult.VmError) > 0 {
			logger.Warnf("eth_estimateGas: EVM execution failed: %v\n", trpcResult.VmError)
			return trpcResult.GasUsed, fmt.Errorf(trpcResult.VmError)
		}
		return trpcResult.GasUsed, nil
	}

	resultIntf, err := common.HandlePandoRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		return "", err
	}

	blockGasLimit := viper.GetUint64(common.CfgPandoBlockGasLimit)
	estimatedGasWithMargin := uint64(1.1 * float64(resultIntf.(tcommon.JSONUint64))) // result should be way below the MAX_UINT_64, so no need to check for overflow
	if estimatedGasWithMargin >= blockGasLimit {
		estimatedGasWithMargin = blockGasLimit
	}
	result = hexutil.EncodeUint64(estimatedGasWithMargin)
	return result, nil
}
