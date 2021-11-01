package ocrkey

import (
	"crypto/ecdsa"

	"github.com/celo-org/celo-blockchain/crypto"
)

type OnChainPublicKey ecdsa.PublicKey

func (k OnChainPublicKey) Address() OnChainSigningAddress {
	return OnChainSigningAddress(crypto.PubkeyToAddress(ecdsa.PublicKey(k)))
}
