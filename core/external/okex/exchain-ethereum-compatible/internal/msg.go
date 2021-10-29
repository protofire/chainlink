package internal

import (
	"math/big"
	"sync/atomic"

	"github.com/celo-org/celo-blockchain/common"
)

// MsgEthereumTx encapsulates an Ethereum transaction as an SDK message.
type MsgEthereumTx struct {
	Data TxData

	// caches
	size atomic.Value
	from atomic.Value
}

func NewMsgEthereumTx(nonce uint64, gasPrice *big.Int, gas uint64, to *common.Address, value *big.Int, input []byte, v, r, s *big.Int) MsgEthereumTx {
	return MsgEthereumTx{
		Data: TxData{
			AccountNonce: nonce,
			Price:        gasPrice,
			GasLimit:     gas,
			Recipient:    to,
			Amount:       value,
			Payload:      input,
			V:            v,
			R:            r,
			S:            s,
		},
	}
}

// TxData implements the Ethereum transaction data structure. It is used
// solely as intended in Ethereum abiding by the protocol.
type TxData struct {
	AccountNonce uint64          `json:"nonce"`
	Price        *big.Int        `json:"gasPrice"`
	GasLimit     uint64          `json:"gas"`
	Recipient    *common.Address `json:"to" rlp:"nil"` // nil means contract creation
	Amount       *big.Int        `json:"value"`
	Payload      []byte          `json:"input"`

	// signature values
	V *big.Int `json:"v"`
	R *big.Int `json:"r"`
	S *big.Int `json:"s"`

	// hash is only used when marshaling to JSON
	Hash *common.Hash `json:"hash" rlp:"-"`
}

// MarshalAmino defines custom encoding scheme for TxData
func (td TxData) MarshalAmino() ([]byte, error) {
	gasPrice, err := MarshalBigInt(td.Price)
	if err != nil {
		return nil, err
	}

	amount, err := MarshalBigInt(td.Amount)
	if err != nil {
		return nil, err
	}

	v, err := MarshalBigInt(td.V)
	if err != nil {
		return nil, err
	}

	r, err := MarshalBigInt(td.R)
	if err != nil {
		return nil, err
	}

	s, err := MarshalBigInt(td.S)
	if err != nil {
		return nil, err
	}

	e := encodableTxData{
		AccountNonce: td.AccountNonce,
		Price:        gasPrice,
		GasLimit:     td.GasLimit,
		Recipient:    td.Recipient,
		Amount:       amount,
		Payload:      td.Payload,
		V:            v,
		R:            r,
		S:            s,
		Hash:         td.Hash,
	}

	return moduleCdc.MarshalBinaryBare(e)
}

// UnmarshalAmino defines custom decoding scheme for TxData
func (td *TxData) UnmarshalAmino(data []byte) error {
	var e encodableTxData
	err := moduleCdc.UnmarshalBinaryBare(data, &e)
	if err != nil {
		return err
	}

	td.AccountNonce = e.AccountNonce
	td.GasLimit = e.GasLimit
	td.Recipient = e.Recipient
	td.Payload = e.Payload
	td.Hash = e.Hash

	price, err := UnmarshalBigInt(e.Price)
	if err != nil {
		return err
	}

	if td.Price != nil {
		td.Price.Set(price)
	} else {
		td.Price = price
	}

	amt, err := UnmarshalBigInt(e.Amount)
	if err != nil {
		return err
	}

	if td.Amount != nil {
		td.Amount.Set(amt)
	} else {
		td.Amount = amt
	}

	v, err := UnmarshalBigInt(e.V)
	if err != nil {
		return err
	}

	if td.V != nil {
		td.V.Set(v)
	} else {
		td.V = v
	}

	r, err := UnmarshalBigInt(e.R)
	if err != nil {
		return err
	}

	if td.R != nil {
		td.R.Set(r)
	} else {
		td.R = r
	}

	s, err := UnmarshalBigInt(e.S)
	if err != nil {
		return err
	}

	if td.S != nil {
		td.S.Set(s)
	} else {
		td.S = s
	}

	return nil
}

type encodableTxData struct {
	AccountNonce uint64          `json:"nonce"`
	Price        string          `json:"gasPrice"`
	GasLimit     uint64          `json:"gas"`
	Recipient    *common.Address `json:"to" rlp:"nil"` // nil means contract creation
	Amount       string          `json:"value"`
	Payload      []byte          `json:"input"`

	// signature values
	V string `json:"v"`
	R string `json:"r"`
	S string `json:"s"`

	// hash is only used when marshaling to JSON
	Hash *common.Hash `json:"hash" rlp:"-"`
}
