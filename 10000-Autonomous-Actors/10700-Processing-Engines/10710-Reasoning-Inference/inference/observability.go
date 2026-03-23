package inference

import (
	"context"
	"log/slog"

	observabilityv1 "olympus.fleet/00SDLC/OlympusGCP/gen/observability"
	"connectrpc.com/connect"
)

type ObservabilityServer struct{}

func (s *ObservabilityServer) WriteLog(ctx context.Context, req *connect.Request[observabilityv1.WriteLogRequest]) (*connect.Response[observabilityv1.WriteLogResponse], error) {
	slog.Info("WriteLog", "log", req.Msg.LogName, "severity", req.Msg.Severity, "msg", req.Msg.Message)
	return connect.NewResponse(&observabilityv1.WriteLogResponse{}), nil
}

func (s *ObservabilityServer) RecordMetric(ctx context.Context, req *connect.Request[observabilityv1.RecordMetricRequest]) (*connect.Response[observabilityv1.RecordMetricResponse], error) {
	slog.Info("RecordMetric", "type", req.Msg.MetricType, "value", req.Msg.Value)
	return connect.NewResponse(&observabilityv1.RecordMetricResponse{}), nil
}

func (s *ObservabilityServer) StartSpan(ctx context.Context, req *connect.Request[observabilityv1.StartSpanRequest]) (*connect.Response[observabilityv1.StartSpanResponse], error) {
	slog.Info("StartSpan", "name", req.Msg.Name, "trace", req.Msg.TraceId)
	return connect.NewResponse(&observabilityv1.StartSpanResponse{SpanId: "span-123"}), nil
}

func (s *ObservabilityServer) EndSpan(ctx context.Context, req *connect.Request[observabilityv1.EndSpanRequest]) (*connect.Response[observabilityv1.EndSpanResponse], error) {
	slog.Info("EndSpan", "span", req.Msg.SpanId)
	return connect.NewResponse(&observabilityv1.EndSpanResponse{}), nil
}
