package ethrpc

import (
	"context"

	"github.com/pandotoken/pando-eth-rpc-adaptor/common"
	trpc "github.com/pandotoken/pando/rpc"

	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_getTransactionByBlockNumberAndIndex -----------------------------------
func (e *EthRPCService) GetTransactionByBlockNumberAndIndex(ctx context.Context, numberStr string, txIndexStr string) (result common.EthGetTransactionResult, err error) {
	logger.Infof("GetTransactionByBlockNumberAndIndex called")
	height := common.GetHeightByTag(numberStr)
	txIndex := common.GetHeightByTag(txIndexStr) //TODO: use common
	client := rpcc.NewRPCClient(common.GetPandoRPCEndpoint())
	rpcRes, rpcErr := client.Call("pando.GetBlockByHeight", trpc.GetBlockByHeightArgs{Height: height})
	return GetIndexedTransactionFromBlock(rpcRes, rpcErr, txIndex)
}
