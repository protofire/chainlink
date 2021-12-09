package utils

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/smartcontractkit/chainlink/core/external/okex/exchain-ethereum-compatible/internal"
)

func Hash(signtx *types.Transaction) (common.Hash, error) {
	// TODO koteld: NOTE
	txSignatures := signtx.RawSignatureValues()[0]
	msg := internal.NewMsgEthereumTx(
		signtx.Nonce(),
		signtx.GasPrice(),
		signtx.Gas(),
		signtx.To(),
		signtx.Value(),
		signtx.Data(),
		txSignatures.V,
		txSignatures.R,
		txSignatures.S,
	)

	bins, err := marshal(msg)
	if err != nil {
		return common.Hash{}, errors.New(fmt.Sprintf("failed to marshal msg: %v", err))
	}

	hash := sha256.Sum256(bins)
	return common.BytesToHash(hash[:]), nil
}

func marshal(msg internal.MsgEthereumTx) ([]byte, error) {
	cdc := internal.GetModuleCdc()
	return cdc.MarshalBinaryLengthPrefixed(msg)
}
