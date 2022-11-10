package ethrpc

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/big"
	"strings"

	"github.com/spf13/viper"
	"github.com/pandoprojects/pando-eth-rpc-adaptor/common"
	tcommon "github.com/pandoprojects/pando/common"
	"github.com/pandoprojects/pando/common/hexutil"
	tcrypto "github.com/pandoprojects/pando/crypto"
	"github.com/pandoprojects/pando/ledger/types"

	trpc "github.com/pandoprojects/pando/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_getBlockByHash -----------------------------------
func (e *EthRPCService) GetBlockByHash(ctx context.Context, hashStr string, txDetails bool) (result common.EthGetBlockResult, err error) {
	logger.Infof("eth_getBlockByHash called, blockHash: %v", hashStr)

	chainIDStr, err := e.ChainId(ctx)
	if err != nil {
		logger.Errorf("Failed to get chainID\n")
		return result, nil
	}
	chainID := new(big.Int)
	chainID.SetString(chainIDStr, 16)

	client := rpcc.NewRPCClient(common.GetPandoRPCEndpoint())
	rpcRes, rpcErr := client.Call("pando.GetBlock", trpc.GetBlockArgs{Hash: tcommon.HexToHash(hashStr)})
	if rpcErr != nil {
		logger.Errorf("eth_getBlockByHash, error: %v", rpcErr)
	}
	return GetBlockFromTRPCResult(chainID, rpcRes, rpcErr, txDetails)
}

func GetBlockFromTRPCResult(chainID *big.Int, rpcRes *rpcc.RPCResponse, rpcErr error, txDetails bool) (result common.EthGetBlockResult, err error) {
	result = common.EthGetBlockResult{}
	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := common.PandoGetBlockResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		if trpcResult.PandoGetBlockResultInner == nil {
			return result, errors.New("empty block")
		}
		result.Transactions = make([]interface{}, 0)
		if txDetails {
			var objmap map[string]json.RawMessage
			json.Unmarshal(jsonBytes, &objmap)
			if objmap["transactions"] != nil {
				var txmaps []map[string]json.RawMessage
				json.Unmarshal(objmap["transactions"], &txmaps)
				for i, omap := range txmaps {
					if types.TxType(trpcResult.Txs[i].Type) == types.TxSmartContract {
						scTx := types.SmartContractTx{}
						json.Unmarshal(omap["raw"], &scTx)

						var ethTx common.EthGetTransactionResult

						ethTx.BlockHash = trpcResult.Hash
						ethTx.BlockHeight = hexutil.Uint64(trpcResult.Height)

						ethTx.From = scTx.From.Address
						if (scTx.To.Address == tcommon.Address{}) {
							ethTx.To = nil // conform to ETH standard
						} else {
							ethTx.To = &scTx.To.Address
						}
						ethTx.GasPrice = hexutil.Uint64(scTx.GasPrice.Uint64())
						ethTx.Gas = hexutil.Uint64(scTx.GasLimit)
						ethTx.Value = hexutil.Uint64(scTx.From.Coins.PTXWei.Uint64())
						ethTx.Input = "0x" + hex.EncodeToString(scTx.Data)
						sigData := scTx.From.Signature.ToBytes()
						ethTx.Nonce = hexutil.Uint64(scTx.From.Sequence) - 1 // off-by-one: Ethereum's account nonce starts from 0, while Pando's account sequnce starts from 1
						//ethTx.TxHash = GetEthTxHash(chainID, ethTx)

						txBytes, _ := types.TxToBytes(&scTx)
						ethTx.TxHash = tcrypto.Keccak256Hash(txBytes)

						GetRSVfromSignature(sigData, &ethTx)

						result.Transactions = append(result.Transactions, ethTx)
						result.GasUsed = hexutil.Uint64(trpcResult.Txs[i].Receipt.GasUsed)
					}
				}
			}
		}
		return trpcResult, nil
	}
	resultIntf, err := common.HandlePandoRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		return result, err
	}
	pando_GetBlockResult := resultIntf.(common.PandoGetBlockResult)
	result.Height = hexutil.Uint64(pando_GetBlockResult.Height)
	result.Hash = pando_GetBlockResult.Hash
	result.Parent = pando_GetBlockResult.Parent
	result.Timestamp = hexutil.Uint64(pando_GetBlockResult.Timestamp.ToInt().Uint64())
	result.Proposer = pando_GetBlockResult.Proposer
	result.TxHash = pando_GetBlockResult.TxHash
	result.StateHash = pando_GetBlockResult.StateHash
	result.GasLimit = hexutil.Uint64(viper.GetUint64(common.CfgPandoBlockGasLimit))
	result.Size = 1000

	for _, tx := range pando_GetBlockResult.Txs {
		if !txDetails && types.TxType(tx.Type) == types.TxSmartContract {
			result.Transactions = append(result.Transactions, tx.Hash)
		}
	}

	result.LogsBloom = "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
	result.ExtraData = "0x"
	result.Nonce = "0x0000000000000000"
	result.Uncles = []tcommon.Hash{}

	return result, nil
}

func GetEthTxHash(chainID *big.Int, ethTx common.EthGetTransactionResult) tcommon.Hash {
	ethTxDataStr := ethTx.Input
	if strings.HasPrefix(ethTx.Input, "0x") {
		ethTxDataStr = ethTxDataStr[2:]
	}
	ethTxData, _ := hex.DecodeString(ethTxDataStr)

	ethTxHash := types.RLPHash([]interface{}{
		ethTx.Nonce,
		new(big.Int).SetUint64(uint64(ethTx.GasPrice)),
		uint64(ethTx.Gas),
		ethTx.To,
		new(big.Int).SetUint64(uint64(ethTx.Value)),
		ethTxData,
		chainID, uint(0), uint(0),
	})
	return ethTxHash
}

// func GetBlockFromTRPCResult(rpcRes *rpcc.RPCResponse, rpcErr error, txDetails bool) (result common.EthGetBlockResult, err error) {
// 	result = common.EthGetBlockResult{}
// 	parse := func(jsonBytes []byte) (interface{}, error) {
// 		trpcResult := trpc.GetBlockResult{}
// 		json.Unmarshal(jsonBytes, &trpcResult)
// 		if trpcResult.GetBlockResultInner == nil {
// 			return result, errors.New("empty block")
// 		}
// 		//result.Transactions = make([]interface{}, len(trpcResult.Txs))
// 		result.Transactions = make([]interface{}, 0)
// 		if txDetails {
// 			var objmap map[string]json.RawMessage
// 			json.Unmarshal(jsonBytes, &objmap)
// 			if objmap["transactions"] != nil {
// 				var txmaps []map[string]json.RawMessage
// 				json.Unmarshal(objmap["transactions"], &txmaps)
// 				for i, omap := range txmaps {
// 					//tx := common.EthGetTransactionResult{}
// 					if types.TxType(trpcResult.Txs[i].Type) == types.TxSmartContract {
// 						scTx := types.SmartContractTx{}
// 						json.Unmarshal(omap["raw"], &scTx)
// 						result.Transactions = append(result.Transactions, scTx)
// 						result.GasUsed = hexutil.Uint64(trpcResult.Txs[i].Receipt.GasUsed)
// 					} else if types.TxType(trpcResult.Txs[i].Type) == types.TxSend {
// 						continue // skip coinbase tx

// 						// sTx := types.SendTx{}
// 						// json.Unmarshal(omap["raw"], &sTx)
// 						// result.Transactions[i] = sTx
// 					} else if types.TxType(trpcResult.Txs[i].Type) == types.TxCoinbase {
// 						continue // skip coinbase tx

// 						// cTx := types.CoinbaseTx{}
// 						// json.Unmarshal(omap["raw"], &cTx)
// 						// tx.From = cTx.Proposer.Address
// 						// tx.Gas = hexutil.Uint64(0)
// 						// tx.Value = hexutil.Uint64(cTx.Proposer.Coins.PTXWei.Uint64())
// 						// tx.Input = "0x"
// 						// data := cTx.Proposer.Signature.ToBytes()
// 						// GetRSVfromSignature(data, &tx)
// 						// result.Transactions[i] = tx
// 					}
// 				}
// 			}
// 		}
// 		return trpcResult, nil
// 	}
// 	resultIntf, err := common.HandlePandoRPCResponse(rpcRes, rpcErr, parse)
// 	if err != nil {
// 		return result, err
// 	}
// 	pando_GetBlockResult := resultIntf.(trpc.GetBlockResult)
// 	result.Height = hexutil.Uint64(pando_GetBlockResult.Height)
// 	result.Hash = pando_GetBlockResult.Hash
// 	result.Parent = pando_GetBlockResult.Parent
// 	result.Timestamp = hexutil.Uint64(pando_GetBlockResult.Timestamp.ToInt().Uint64())
// 	result.Proposer = pando_GetBlockResult.Proposer
// 	result.TxHash = pando_GetBlockResult.TxHash
// 	result.StateHash = pando_GetBlockResult.StateHash
// 	// for i, tx := range pando_GetBlockResult.Txs {
// 	// 	if txDetails && (types.TxType(tx.Type) == types.TxSmartContract || types.TxType(tx.Type) == types.TxSend || types.TxType(tx.Type) == types.TxCoinbase) {
// 	// 		//already handled
// 	// 	} else {
// 	// 		result.Transactions[i] = tx.Hash
// 	// 	}
// 	// }
// 	result.GasLimit = hexutil.Uint64(viper.GetUint64(common.CfgPandoBlockGasLimit))

// 	result.LogsBloom = "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
// 	result.ExtraData = "0x"
// 	result.Nonce = "0x0000000000000000"
// 	result.Uncles = []tcommon.Hash{}

// 	return result, nil
// }
