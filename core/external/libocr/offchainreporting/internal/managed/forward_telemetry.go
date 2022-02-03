package managed

import (
	"context"

	"github.com/smartcontractkit/chainlink/core/external/libocr/offchainreporting/internal/serialization/protobuf"
	"github.com/smartcontractkit/chainlink/core/external/libocr/offchainreporting/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting/types"
	"google.golang.org/protobuf/proto"
)

// forwardTelemetry receives monitoring events on chTelemetry, serializes them, and forwards
// them to monitoringEndpoint
func forwardTelemetry(
	ctx context.Context,

	logger loghelper.LoggerWithContext,
	monitoringEndpoint types.MonitoringEndpoint,

	chTelemetry <-chan *protobuf.TelemetryWrapper,
) {
	for {
		select {
		case t, ok := <-chTelemetry:
			if !ok {
				// This isn't supposed to happen, but we still handle this case gracefully,
				// just in case...
				logger.Error("forwardTelemetry: chTelemetry closed unexpectedly. exiting", nil)
				return
			}
			bin, err := proto.Marshal(t)
			if err != nil {
				logger.Error("forwardTelemetry: failed to Marshal protobuf", types.LogFields{
					"proto": t,
					"error": err,
				})
				break
			}
			if monitoringEndpoint != nil {
				monitoringEndpoint.SendLog(bin)
			}
		case <-ctx.Done():
			logger.Info("forwardTelemetry: exiting", nil)
			return
		}
	}
}
