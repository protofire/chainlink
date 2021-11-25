package fluxmonitorv2

import (
	"github.com/klaytn/klaytn/common"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
)

//go:generate mockery --name KeyStoreInterface --output ./mocks/ --case=underscore

// KeyStoreInterface defines an interface to interact with the keystore
type KeyStoreInterface interface {
	SendingKeys() ([]ethkey.KeyV2, error)
	GetRoundRobinAddress(...common.Address) (common.Address, error)
}

// KeyStore implements KeyStoreInterface
type KeyStore struct {
	keystore.Eth
}

// NewKeyStore initializes a new keystore
func NewKeyStore(ks keystore.Eth) *KeyStore {
	return &KeyStore{ks}
}
