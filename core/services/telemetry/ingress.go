package telemetry

import (
	"context"

	"github.com/klaytn/klaytn/common"
	ocrtypes "github.com/smartcontractkit/chainlink/core/external/libocr/offchainreporting/types"
	"github.com/smartcontractkit/chainlink/core/services/synchronization"
)

type IngressAgentWrapper struct {
	telemetryIngressClient synchronization.TelemetryIngressClient
}

func NewIngressAgentWrapper(telemetryIngressClient synchronization.TelemetryIngressClient) *IngressAgentWrapper {
	return &IngressAgentWrapper{telemetryIngressClient}
}

func (t *IngressAgentWrapper) GenMonitoringEndpoint(addr common.Address) ocrtypes.MonitoringEndpoint {
	return NewIngressAgent(t.telemetryIngressClient, addr)
}

type IngressAgent struct {
	telemetryIngressClient synchronization.TelemetryIngressClient
	contractAddress        common.Address
}

func NewIngressAgent(telemetryIngressClient synchronization.TelemetryIngressClient, contractAddress common.Address) *IngressAgent {
	return &IngressAgent{
		telemetryIngressClient,
		contractAddress,
	}
}

// SendLog sends a telemetry log to the ingress server
func (t *IngressAgent) SendLog(telemetry []byte) {
	payload := synchronization.TelemPayload{
		Ctx:             context.Background(),
		Telemetry:       telemetry,
		ContractAddress: t.contractAddress,
	}
	t.telemetryIngressClient.Send(payload)
}
