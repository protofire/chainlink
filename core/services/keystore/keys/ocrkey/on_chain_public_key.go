package ocrkey

import (
	"crypto/ecdsa"

	"github.com/klaytn/klaytn/crypto"
)

type OnChainPublicKey ecdsa.PublicKey

func (k OnChainPublicKey) Address() OnChainSigningAddress {
	return OnChainSigningAddress(crypto.PubkeyToAddress(ecdsa.PublicKey(k)))
}
