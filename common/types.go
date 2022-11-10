package common

import (
	"github.com/pandotoken/pando/blockchain"
	"github.com/pandotoken/pando/common"
	tcommon "github.com/pandotoken/pando/common"
	"github.com/pandotoken/pando/common/hexutil"
	tcore "github.com/pandotoken/pando/core"
	"github.com/pandotoken/pando/ledger/types"
	trpc "github.com/pandotoken/pando/rpc"
)

type Bytes8 [8]byte

// type Bytes []byte
type EthGetTransactionResult struct {
	BlockHash        tcommon.Hash     `json:"blockHash"`
	BlockHeight      hexutil.Uint64   `json:"blockNumber"`
	From             tcommon.Address  `json:"from"`
	To               *tcommon.Address `json:"to"`
	Gas              hexutil.Uint64   `json:"gas"`
	GasPrice         hexutil.Uint64   `json:"gasPrice"`
	TxHash           tcommon.Hash     `json:"hash"`
	Nonce            hexutil.Uint64   `json:"nonce"`
	Input            string           `json:"input"`
	TransactionIndex hexutil.Uint64   `json:"transactionIndex"`
	Value            hexutil.Uint64   `json:"value"`
	V                hexutil.Uint64   `json:"v"` //ECDSA recovery id
	R                tcommon.Hash     `json:"r"` //ECDSA signature r
	S                tcommon.Hash     `json:"s"` //ECDSA signature s
}

type EthGetBlockResult struct {
	Height    hexutil.Uint64  `json:"number"`
	Hash      tcommon.Hash    `json:"hash"`
	Parent    tcommon.Hash    `json:"parentHash"`
	Timestamp hexutil.Uint64  `json:"timestamp"`
	Proposer  tcommon.Address `json:"miner"`
	TxHash    tcommon.Hash    `json:"transactionsRoot"`
	StateHash tcommon.Hash    `json:"stateRoot"`

	ReiceptHash     tcommon.Hash   `json:"receiptsRoot"`
	Nonce           string         `json:"nonce"`
	Sha3Uncles      tcommon.Hash   `json:"sha3Uncles"`
	LogsBloom       string         `json:"logsBloom"`
	Difficulty      hexutil.Uint64 `json:"difficulty"`
	TotalDifficulty hexutil.Uint64 `json:"totalDifficulty"`
	Size            hexutil.Uint64 `json:"size"`
	GasLimit        hexutil.Uint64 `json:"gasLimit"`
	GasUsed         hexutil.Uint64 `json:"gasUsed"`
	ExtraData       string         `json:"extraData"`
	Uncles          []tcommon.Hash `json:"uncles"`
	Transactions    []interface{}  `json:"transactions"`
}

type EthSyncingResult struct {
	StartingBlock hexutil.Uint64 `json:"startingBlock"`
	CurrentBlock  hexutil.Uint64 `json:"currentBlock"`
	HighestBlock  hexutil.Uint64 `json:"highestBlock"`
	PulledStates  hexutil.Uint64 `json:"pulledStates"` //pulledStates is the number it already downloaded
	KnownStates   hexutil.Uint64 `json:"knownStates"`  //knownStates is the number of trie nodes that the sync algo knows about
}

type EthGetReceiptResult struct {
	BlockHash         tcommon.Hash    `json:"blockHash"`
	BlockHeight       hexutil.Uint64  `json:"blockNumber"`
	TxHash            tcommon.Hash    `json:"transactionHash"`
	TransactionIndex  hexutil.Uint64  `json:"transactionIndex"`
	ContractAddress   tcommon.Address `json:"contractAddress"`
	From              tcommon.Address `json:"from"`
	To                tcommon.Address `json:"to"`
	GasUsed           hexutil.Uint64  `json:"gasUsed"`
	CumulativeGasUsed hexutil.Uint64  `json:"cumulativeGasUsed"`
	Logs              []EthLogObj     `json:"logs"`
	LogsBloom         string          `json:"logsBloom"`
	Status            hexutil.Uint64  `json:"status"`
}

type Tx struct {
	types.Tx `json:"raw"`
	Type     byte                       `json:"type"`
	Hash     tcommon.Hash               `json:"hash"`
	Receipt  *blockchain.TxReceiptEntry `json:"receipt"`
}

type EthLogObj struct {
	Address          tcommon.Address `json:"address"`
	BlockHash        tcommon.Hash    `json:"blockHash"`
	BlockHeight      hexutil.Uint64  `json:"blockNumber"`
	LogIndex         hexutil.Uint64  `json:"logIndex"`
	Topics           []tcommon.Hash  `json:"topics"`
	TxHash           tcommon.Hash    `json:"transactionHash"`
	TransactionIndex hexutil.Uint64  `json:"transactionIndex"`
	Data             string          `json:"data"`
	Type             string          `json:"type"`
	//Removed          bool            `json:"removed"`

}

type EthSmartContractArgObj struct {
	From     tcommon.Address `json:"from"`
	To       tcommon.Address `json:"to"`
	Gas      string          `json:"gas"`
	GasPrice string          `json:"gasPrice"`
	Value    string          `json:"value"`
	Data     string          `json:"data"`
}

type PandoGetBlockResult struct {
	*PandoGetBlockResultInner
}
type PandoGetBlocksResult []*PandoGetBlockResultInner

type PandoGetBlockResultInner struct {
	ChainID            string                    `json:"chain_id"`
	Epoch              tcommon.JSONUint64        `json:"epoch"`
	Height             tcommon.JSONUint64        `json:"height"`
	Parent             tcommon.Hash              `json:"parent"`
	TxHash             tcommon.Hash              `json:"transactions_hash"`
	StateHash          tcommon.Hash              `json:"state_hash"`
	Timestamp          *tcommon.JSONBig          `json:"timestamp"`
	Proposer           tcommon.Address           `json:"proposer"`
	HCC                tcore.CommitCertificate   `json:"hcc"`
	GuardianVotes      *tcore.AggregatedVotes    `json:"guardian_votes"`
	RametronenterpriseVotes *tcore.AggregatedRametronenterpriseVotes `json:"rametronenterprise_votes"`

	Children []common.Hash     `json:"children"`
	Status   tcore.BlockStatus `json:"status"`

	Hash common.Hash `json:"hash"`
	Txs  []trpc.Tx   `json:"transactions"`
}
