package protocol

import "github.com/smartcontractkit/libocr/offchainreporting/types"

type TelemetrySender interface {
	RoundStarted(
		configDigest types.ConfigDigest,
		epoch uint32,
		round uint8,
		leader types.OracleID,
	)
}
