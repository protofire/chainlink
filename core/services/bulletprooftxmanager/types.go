package bulletprooftxmanager

import (
	"encoding/json"
	gethTypes "github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// Receipt represents an ethereum receipt.
//
// Copied from go-ethereum: https://github.com/ethereum/go-ethereum/blob/ce9a289fa48e0d2593c4aaa7e207c8a5dd3eaa8a/core/types/receipt.go#L50
//
// We use our own version because Geth's version specifies various
// gencodec:"required" fields which cause unhelpful errors when unmarshalling
// from an empty JSON object which can happen in the batch fetcher.

// TODO koteld: NOTE removed all the fields not supported by Klaytn Receipt type
type Receipt struct {
	Status          uint64          `json:"status"`
	Bloom           gethTypes.Bloom `json:"logsBloom"`
	Logs            []*Log          `json:"logs"`
	TxHash          common.Hash     `json:"transactionHash"`
	ContractAddress common.Address  `json:"contractAddress"`
	GasUsed         uint64          `json:"gasUsed"`
}

// FromGethReceipt converts a gethTypes.Receipt to a Receipt
func FromGethReceipt(gr *gethTypes.Receipt) *Receipt {
	if gr == nil {
		return nil
	}
	logs := make([]*Log, len(gr.Logs))
	for i, glog := range gr.Logs {
		logs[i] = FromGethLog(glog)
	}
	return &Receipt{
		uint64(gr.Status),
		gr.Bloom,
		logs,
		gr.TxHash,
		gr.ContractAddress,
		gr.GasUsed,
	}
}

// IsZero returns true if receipt is the zero receipt
// Batch calls to the RPC will return a pointer to an empty Receipt struct
// Easiest way to check if the receipt was missing is to see if the hash is 0x0
// Real receipts will always have the TxHash set
func (r Receipt) IsZero() bool {
	return r.TxHash == utils.EmptyHash
}

// MarshalJSON marshals Receipt as JSON.
// Copied from: https://github.com/ethereum/go-ethereum/blob/ce9a289fa48e0d2593c4aaa7e207c8a5dd3eaa8a/core/types/gen_receipt_json.go
func (r Receipt) MarshalJSON() ([]byte, error) {
	type Receipt struct {
		Status          hexutil.Uint64  `json:"status"`
		Bloom           gethTypes.Bloom `json:"logsBloom"`
		Logs            []*Log          `json:"logs"`
		TxHash          common.Hash     `json:"transactionHash"`
		ContractAddress common.Address  `json:"contractAddress"`
		GasUsed         hexutil.Uint64  `json:"gasUsed"`
	}
	var enc Receipt
	enc.Status = hexutil.Uint64(r.Status)
	enc.Bloom = r.Bloom
	enc.Logs = r.Logs
	enc.TxHash = r.TxHash
	enc.ContractAddress = r.ContractAddress
	enc.GasUsed = hexutil.Uint64(r.GasUsed)
	return json.Marshal(&enc)
}

// UnmarshalJSON unmarshals from JSON.
func (r *Receipt) UnmarshalJSON(input []byte) error {
	type Receipt struct {
		Status          *hexutil.Uint64  `json:"status"`
		Bloom           *gethTypes.Bloom `json:"logsBloom"`
		Logs            []*Log           `json:"logs"`
		TxHash          *common.Hash     `json:"transactionHash"`
		ContractAddress *common.Address  `json:"contractAddress"`
		GasUsed         *hexutil.Uint64  `json:"gasUsed"`
	}
	var dec Receipt
	if err := json.Unmarshal(input, &dec); err != nil {
		return errors.Wrap(err, "could not unmarshal receipt")
	}
	if dec.Status != nil {
		r.Status = uint64(*dec.Status)
	}
	//if dec.CumulativeGasUsed != nil {
	//	r.CumulativeGasUsed = uint64(*dec.CumulativeGasUsed)
	//}
	if dec.Bloom != nil {
		r.Bloom = *dec.Bloom
	}
	r.Logs = dec.Logs
	if dec.TxHash != nil {
		r.TxHash = *dec.TxHash
	}
	if dec.ContractAddress != nil {
		r.ContractAddress = *dec.ContractAddress
	}
	if dec.GasUsed != nil {
		r.GasUsed = uint64(*dec.GasUsed)
	}
	return nil
}

// Log represents a contract log event.
//
// Copied from go-ethereum: https://github.com/ethereum/go-ethereum/blob/ce9a289fa48e0d2593c4aaa7e207c8a5dd3eaa8a/core/types/log.go
//
// We use our own version because Geth's version specifies various
// gencodec:"required" fields which cause unhelpful errors when unmarshalling
// from an empty JSON object which can happen in the batch fetcher.
type Log struct {
	Address     common.Address `json:"address"`
	Topics      []common.Hash  `json:"topics"`
	Data        []byte         `json:"data"`
	BlockNumber uint64         `json:"blockNumber"`
	TxHash      common.Hash    `json:"transactionHash"`
	TxIndex     uint           `json:"transactionIndex"`
	BlockHash   common.Hash    `json:"blockHash"`
	Index       uint           `json:"logIndex"`
	Removed     bool           `json:"removed"`
}

// FromGethLog converts a gethTypes.Log to a Log
func FromGethLog(gl *gethTypes.Log) *Log {
	if gl == nil {
		return nil
	}
	return &Log{
		gl.Address,
		gl.Topics,
		gl.Data,
		gl.BlockNumber,
		gl.TxHash,
		gl.TxIndex,
		gl.BlockHash,
		gl.Index,
		gl.Removed,
	}
}

// MarshalJSON marshals as JSON.
func (l Log) MarshalJSON() ([]byte, error) {
	type Log struct {
		Address     common.Address `json:"address"`
		Topics      []common.Hash  `json:"topics"`
		Data        hexutil.Bytes  `json:"data"`
		BlockNumber hexutil.Uint64 `json:"blockNumber"`
		TxHash      common.Hash    `json:"transactionHash"`
		TxIndex     hexutil.Uint   `json:"transactionIndex"`
		BlockHash   common.Hash    `json:"blockHash"`
		Index       hexutil.Uint   `json:"logIndex"`
		Removed     bool           `json:"removed"`
	}
	var enc Log
	enc.Address = l.Address
	enc.Topics = l.Topics
	enc.Data = l.Data
	enc.BlockNumber = hexutil.Uint64(l.BlockNumber)
	enc.TxHash = l.TxHash
	enc.TxIndex = hexutil.Uint(l.TxIndex)
	enc.BlockHash = l.BlockHash
	enc.Index = hexutil.Uint(l.Index)
	enc.Removed = l.Removed
	return json.Marshal(&enc)
}

// UnmarshalJSON unmarshals from JSON.
func (l *Log) UnmarshalJSON(input []byte) error {
	type Log struct {
		Address     *common.Address `json:"address"`
		Topics      []common.Hash   `json:"topics"`
		Data        *hexutil.Bytes  `json:"data"`
		BlockNumber *hexutil.Uint64 `json:"blockNumber"`
		TxHash      *common.Hash    `json:"transactionHash"`
		TxIndex     *hexutil.Uint   `json:"transactionIndex"`
		BlockHash   *common.Hash    `json:"blockHash"`
		Index       *hexutil.Uint   `json:"logIndex"`
		Removed     *bool           `json:"removed"`
	}
	var dec Log
	if err := json.Unmarshal(input, &dec); err != nil {
		return errors.Wrap(err, "coult not unmarshal log")
	}
	if dec.Address != nil {
		l.Address = *dec.Address
	}
	l.Topics = dec.Topics
	if dec.Data != nil {
		l.Data = *dec.Data
	}
	if dec.BlockNumber != nil {
		l.BlockNumber = uint64(*dec.BlockNumber)
	}
	if dec.TxHash != nil {
		l.TxHash = *dec.TxHash
	}
	if dec.TxIndex != nil {
		l.TxIndex = uint(*dec.TxIndex)
	}
	if dec.BlockHash != nil {
		l.BlockHash = *dec.BlockHash
	}
	if dec.Index != nil {
		l.Index = uint(*dec.Index)
	}
	if dec.Removed != nil {
		l.Removed = *dec.Removed
	}
	return nil
}
