package workflows

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/metrics"

	localMonitoring "github.com/smartcontractkit/chainlink/v2/core/monitoring"
)

var registerTriggerFailureCounter metric.Int64Counter
var workflowsRunningGauge metric.Int64Gauge
var capabilityInvocationCounter metric.Int64Counter
var workflowExecutionLatencyGauge metric.Int64Gauge // ms
var workflowStepErrorCounter metric.Int64Counter
var engineHeartbeatCounter metric.Int64UpDownCounter

func initMonitoringResources() (err error) {
	registerTriggerFailureCounter, err = beholder.GetMeter().Int64Counter("platform.engine.register_trigger.failures")
	if err != nil {
		return fmt.Errorf("failed to register trigger failure counter: %w", err)
	}

	workflowsRunningGauge, err = beholder.GetMeter().Int64Gauge("platform.engine.workflows.count")
	if err != nil {
		return fmt.Errorf("failed to register workflows running gauge: %w", err)
	}

	capabilityInvocationCounter, err = beholder.GetMeter().Int64Counter("platform.engine.capabilities_invoked.count")
	if err != nil {
		return fmt.Errorf("failed to register capability invocation counter: %w", err)
	}

	workflowExecutionLatencyGauge, err = beholder.GetMeter().Int64Gauge("platform.engine.workflow.time")
	if err != nil {
		return fmt.Errorf("failed to register workflow execution latency gauge: %w", err)
	}

	workflowStepErrorCounter, err = beholder.GetMeter().Int64Counter("platform.engine.workflow.errors")
	if err != nil {
		return fmt.Errorf("failed to register workflow step error counter: %w", err)
	}

	engineHeartbeatCounter, err = beholder.GetMeter().Int64UpDownCounter("platform.engine.heartbeat")
	if err != nil {
		return fmt.Errorf("failed to register engine heartbeat counter: %w", err)
	}

	return nil
}

// workflowsMetricLabeler wraps monitoring.MetricsLabeler to provide workflow specific utilities
// for monitoring resources
type workflowsMetricLabeler struct {
	metrics.Labeler
}

func (c workflowsMetricLabeler) with(keyValues ...string) workflowsMetricLabeler {
	return workflowsMetricLabeler{c.With(keyValues...)}
}

func (c workflowsMetricLabeler) incrementRegisterTriggerFailureCounter(ctx context.Context) {
	otelLabels := localMonitoring.KvMapToOtelAttributes(c.Labels)
	registerTriggerFailureCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (c workflowsMetricLabeler) incrementCapabilityInvocationCounter(ctx context.Context) {
	otelLabels := localMonitoring.KvMapToOtelAttributes(c.Labels)
	capabilityInvocationCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (c workflowsMetricLabeler) updateWorkflowExecutionLatencyGauge(ctx context.Context, val int64) {
	otelLabels := localMonitoring.KvMapToOtelAttributes(c.Labels)
	workflowExecutionLatencyGauge.Record(ctx, val, metric.WithAttributes(otelLabels...))
}

func (c workflowsMetricLabeler) incrementTotalWorkflowStepErrorsCounter(ctx context.Context) {
	otelLabels := localMonitoring.KvMapToOtelAttributes(c.Labels)
	workflowStepErrorCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (c workflowsMetricLabeler) updateTotalWorkflowsGauge(ctx context.Context, val int64) {
	otelLabels := localMonitoring.KvMapToOtelAttributes(c.Labels)
	workflowsRunningGauge.Record(ctx, val, metric.WithAttributes(otelLabels...))
}

func (c workflowsMetricLabeler) incrementEngineHeartbeatCounter(ctx context.Context) {
	otelLabels := localMonitoring.KvMapToOtelAttributes(c.Labels)
	engineHeartbeatCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

// Observability keys
const (
	cIDKey  = "capabilityID"
	tIDKey  = "triggerID"
	wIDKey  = "workflowID"
	eIDKey  = "workflowExecutionID"
	wnKey   = "workflowName"
	woIDKey = "workflowOwner"
	sIDKey  = "stepID"
	sRKey   = "stepRef"
)

var orderedLabelKeys = []string{sRKey, sIDKey, tIDKey, cIDKey, eIDKey, wIDKey}
