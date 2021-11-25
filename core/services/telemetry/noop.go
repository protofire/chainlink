package telemetry

import (
	"github.com/klaytn/klaytn/common"
	ocrtypes "github.com/smartcontractkit/chainlink/core/external/libocr/offchainreporting/types"
)

type NoopAgent struct {
}

// SendLog sends a telemetry log to the explorer
func (t *NoopAgent) SendLog(log []byte) {
}

// GenMonitoringEndpoint creates a monitoring endpoint for telemetry
func (t *NoopAgent) GenMonitoringEndpoint(addr common.Address) ocrtypes.MonitoringEndpoint {
	return t
}
