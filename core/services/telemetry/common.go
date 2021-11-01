package telemetry

import (
	"github.com/celo-org/celo-blockchain/common"
	ocrtypes "github.com/smartcontractkit/chainlink/core/external/libocr/offchainreporting/types"
)

type MonitoringEndpointGenerator interface {
	GenMonitoringEndpoint(addr common.Address) ocrtypes.MonitoringEndpoint
}
