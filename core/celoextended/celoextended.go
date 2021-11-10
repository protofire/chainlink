package celoextended

import (
	"errors"
	"strings"
	"sync"

	"github.com/celo-org/celo-blockchain/accounts/abi"
)

// ErrMethodNotSupported is returned when a method is not supported.
var ErrMethodNotSupported = errors.New("method is not supported")

var MinerGasCeil uint64 = 8000000

// MetaData collects all metadata for a bound contract.
type MetaData struct {
	mu   sync.Mutex
	Sigs map[string]string
	Bin  string
	ABI  string
	ab   *abi.ABI
}

func (m *MetaData) GetAbi() (*abi.ABI, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.ab != nil {
		return m.ab, nil
	}
	if parsed, err := abi.JSON(strings.NewReader(m.ABI)); err != nil {
		return nil, err
	} else {
		m.ab = &parsed
	}
	return m.ab, nil
}
