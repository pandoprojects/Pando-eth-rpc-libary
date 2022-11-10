package ethrpc

import (
	"context"
	"fmt"

	"github.com/pandoprojects/pando-eth-rpc-adaptor/common"
	hexutil "github.com/pandoprojects/pando/common/hexutil"
)

// ------------------------------- eth_blockNumber -----------------------------------

func (e *EthRPCService) BlockNumber(ctx context.Context) (result string, err error) {
	logger.Infof("eth_blockNumber called")

	blockNumber, err := common.GetCurrentHeight()

	if err != nil {
		return "", err
	}

	result = hexutil.EncodeUint64(uint64(blockNumber))
	fmt.Printf("eth_blockNumber result:%+v\n", result)
	return result, nil
}
