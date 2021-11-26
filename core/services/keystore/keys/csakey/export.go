package csakey

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/klaytnextended"
	"github.com/smartcontractkit/chainlink/core/utils"
)

const keyTypeIdentifier = "CSA"

func FromEncryptedJSON(keyJSON []byte, password string) (KeyV2, error) {
	var export EncryptedCSAKeyExport
	if err := json.Unmarshal(keyJSON, &export); err != nil {
		return KeyV2{}, err
	}
	privKey, err := klaytnextended.DecryptDataV3(export.Crypto, adulteratedPassword(password))
	if err != nil {
		return KeyV2{}, errors.Wrap(err, "failed to decrypt CSA key")
	}
	key := Raw(privKey).Key()
	return key, nil
}

type EncryptedCSAKeyExport struct {
	KeyType   string                    `json:"keyType"`
	PublicKey string                    `json:"publicKey"`
	Crypto    klaytnextended.CryptoJSON `json:"crypto"`
}

func (key KeyV2) ToEncryptedJSON(password string, scryptParams utils.ScryptParams) (export []byte, err error) {
	cryptoJSON, err := klaytnextended.EncryptDataV3(
		key.Raw(),
		[]byte(adulteratedPassword(password)),
		scryptParams.N,
		scryptParams.P,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "could not encrypt Eth key")
	}
	encryptedOCRKExport := EncryptedCSAKeyExport{
		KeyType:   keyTypeIdentifier,
		PublicKey: key.PublicKeyString(),
		Crypto:    cryptoJSON,
	}
	return json.Marshal(encryptedOCRKExport)
}

func adulteratedPassword(password string) string {
	return "csakey" + password
}
