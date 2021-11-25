package keystest

import (
	"crypto/ecdsa"
	"crypto/rand"

	"github.com/klaytn/klaytn/accounts/keystore"
	"github.com/klaytn/klaytn/crypto"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
)

// NewKey pulled from geth
func NewKey() (key keystore.KeyV3, err error) {
	privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		return key, err
	}

	id := uuid.NewRandom()
	if err != nil {
		return key, errors.Errorf("Could not create random uuid: %v", err)
	}

	return keystore.KeyV3{
		Id:         id,
		Address:    crypto.PubkeyToAddress(privateKeyECDSA.PublicKey),
		PrivateKey: privateKeyECDSA,
	}, nil
}
