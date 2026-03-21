package main

import (
	"context"
	"testing"

	observabilityv1 "olympus.fleet/00SDLC/OlympusGCP-Observability/40000-Communication-Contracts/40400-Protocol-Synthetics/connect-rpc/gen/v1/observability"
	"olympus.fleet/00SDLC/OlympusGCP-Observability/10000-Autonomous-Actors/10700-Processing-Engines/10710-Reasoning-Inference/inference"
	"connectrpc.com/connect"
)

func TestObservabilityServer(t *testing.T) {
	server := &inference.ObservabilityServer{}
	ctx := context.Background()

	// Test WriteLog
	logReq := connect.NewRequest(&observabilityv1.WriteLogRequest{
		LogName:  "syslog",
		Severity: "INFO",
		Message:  "Everything normal",
	})
	_, err := server.WriteLog(ctx, logReq)
	if err != nil {
		t.Fatalf("WriteLog failed: %v", err)
	}

	// Test RecordMetric
	metricReq := connect.NewRequest(&observabilityv1.RecordMetricRequest{
		MetricType: "request_count",
		Value:      1.0,
	})
	_, err = server.RecordMetric(ctx, metricReq)
	if err != nil {
		t.Fatalf("RecordMetric failed: %v", err)
	}

	// Test Spans
	spanReq := connect.NewRequest(&observabilityv1.StartSpanRequest{
		Name:    "db_query",
		TraceId: "trace-123",
	})
	spanRes, err := server.StartSpan(ctx, spanReq)
	if err != nil {
		t.Fatalf("StartSpan failed: %v", err)
	}
	if spanRes.Msg.SpanId == "" {
		t.Error("Expected span ID, got empty string")
	}

	endReq := connect.NewRequest(&observabilityv1.EndSpanRequest{
		SpanId: spanRes.Msg.SpanId,
	})
	_, err = server.EndSpan(ctx, endReq)
	if err != nil {
		t.Fatalf("EndSpan failed: %v", err)
	}
}
