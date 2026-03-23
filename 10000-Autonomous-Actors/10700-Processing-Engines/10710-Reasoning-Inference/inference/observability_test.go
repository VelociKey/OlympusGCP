package inference

import (
	"context"
	"testing"

	observabilityv1 "olympus.fleet/00SDLC/OlympusGCP/gen/observability"
	"connectrpc.com/connect"
)

func TestObservabilityServer_CoverageExpansion(t *testing.T) {
	server := &ObservabilityServer{}
	ctx := context.Background()

	// 1. Test WriteLog
	_, err := server.WriteLog(ctx, connect.NewRequest(&observabilityv1.WriteLogRequest{
		LogName: "app",
		Severity: "INFO",
		Message: "log",
	}))
	if err != nil {
		t.Error("WriteLog failed")
	}

	// 2. Test RecordMetric
	_, err = server.RecordMetric(ctx, connect.NewRequest(&observabilityv1.RecordMetricRequest{
		MetricType: "counter",
		Value: 1.0,
	}))
	if err != nil {
		t.Error("RecordMetric failed")
	}

	// 3. Test StartSpan
	res, err := server.StartSpan(ctx, connect.NewRequest(&observabilityv1.StartSpanRequest{
		Name: "op1",
	}))
	if err != nil || res.Msg.SpanId == "" {
		t.Error("StartSpan failed")
	}

	// 4. Test EndSpan
	_, err = server.EndSpan(ctx, connect.NewRequest(&observabilityv1.EndSpanRequest{
		SpanId: res.Msg.SpanId,
	}))
	if err != nil {
		t.Error("EndSpan failed")
	}
}
