package vrfkey

import (
	"encoding/json"
	"fmt"
	"math/big"

	keystore "github.com/celo-org/celo-blockchain/accounts/keystore"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
	"go.dedis.ch/kyber/v3"
)

// PrivateKey represents the secret used to construct a VRF proof.
//
// Don't serialize directly, use Encrypt method, with user-supplied passphrase.
// The unencrypted PrivateKey struct should only live in-memory.
//
// Only use it if you absolutely need it (i.e., for a novel crypto protocol.)
// Implement whatever cryptography you need on this struct, so your callers
// don't need to know the secret key explicitly. (See, e.g., MarshaledProof.)
type PrivateKey struct {
	k         kyber.Scalar
	PublicKey secp256k1.PublicKey
}

// newPrivateKey(k) is k wrapped in a PrivateKey along with corresponding
// PublicKey, or an error. Internal use only. Use cltest.StoredVRFKey for stable
// testing key, or CreateKey if you don't need determinism.
func newPrivateKey(rawKey *big.Int) (*PrivateKey, error) {
	if rawKey.Cmp(secp256k1.GroupOrder) >= 0 || rawKey.Cmp(big.NewInt(0)) <= 0 {
		return nil, fmt.Errorf("secret key must be in {1, ..., #secp256k1 - 1}")
	}
	sk := &PrivateKey{}
	sk.k = secp256k1.IntToScalar(rawKey)
	pk, err := suite.Point().Mul(sk.k, nil).MarshalBinary()
	if err != nil {
		panic(errors.Wrapf(err, "could not marshal public key"))
	}
	if len(pk) != secp256k1.CompressedPublicKeyLength {
		panic(fmt.Errorf("public key %x has wrong length", pk))
	}
	if l := copy(sk.PublicKey[:], pk); l != secp256k1.CompressedPublicKeyLength {
		panic(fmt.Errorf("failed to copy correct length in serialized public key"))
	}
	return sk, nil
}

func (k PrivateKey) ToV2() KeyV2 {
	return KeyV2{
		k:         &k.k,
		PublicKey: k.PublicKey,
	}
}

// fromGethKey returns the vrfkey representation of gethKey. Do not abuse this
// to convert an ethereum key into a VRF key!
func fromGethKey(gethKey *keystore.Key) *PrivateKey {
	secretKey := secp256k1.IntToScalar(gethKey.PrivateKey.D)
	rawPublicKey, err := secp256k1.ScalarToPublicPoint(secretKey).MarshalBinary()
	if err != nil {
		panic(err) // Only way this can happen is out-of-memory failure
	}
	var publicKey secp256k1.PublicKey
	copy(publicKey[:], rawPublicKey)
	return &PrivateKey{secretKey, publicKey}
}

func (k *PrivateKey) String() string {
	return fmt.Sprintf("PrivateKey{k: <redacted>, PublicKey: %s}", k.PublicKey)
}

// GoString reduces the risk of accidentally logging the private key
func (k *PrivateKey) GoString() string {
	return k.String()
}

// Decrypt returns the PrivateKey in e, decrypted via auth, or an error
func Decrypt(e EncryptedVRFKey, auth string) (*PrivateKey, error) {
	// NOTE: We do this shuffle to an anonymous struct
	// solely to add a a throwaway UUID, so we can leverage
	// the keystore.DecryptKey from the geth which requires it
	// as of 1.10.0.
	keyJSON, err := json.Marshal(struct {
		Address string              `json:"address"`
		Crypto  keystore.CryptoJSON `json:"crypto"`
		Version int                 `json:"version"`
		Id      string              `json:"id"`
	}{
		Address: e.VRFKey.Address,
		Crypto:  e.VRFKey.Crypto,
		Version: e.VRFKey.Version,
		Id:      uuid.New().String(),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "while marshaling key for decryption")
	}
	gethKey, err := keystore.DecryptKey(keyJSON, adulteratedPassword(auth))
	if err != nil {
		return nil, errors.Wrapf(err, "could not decrypt key %s",
			e.PublicKey.String())
	}
	return fromGethKey(gethKey), nil
}
