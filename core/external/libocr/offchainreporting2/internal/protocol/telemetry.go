package protocol

import (
	"github.com/smartcontractkit/chainlink/core/external/libocr/commontypes"
	"github.com/smartcontractkit/chainlink/core/external/libocr/offchainreporting2/types"
)

type TelemetrySender interface {
	RoundStarted(
		configDigest types.ConfigDigest,
		epoch uint32,
		round uint8,
		leader commontypes.OracleID,
	)
}
