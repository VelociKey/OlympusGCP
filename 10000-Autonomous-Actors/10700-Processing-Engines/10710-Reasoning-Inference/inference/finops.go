package inference

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	finopsv1 "olympus.fleet/00SDLC/OlympusGCP/gen/finops"
	computev1 "olympus.fleet/00SDLC/OlympusGCP/gen/compute"
	"olympus.fleet/00SDLC/OlympusGCP/gen/google/cloud/compute/computev1connect"
	"connectrpc.com/connect"
)

type FinOpsServer struct {
	computeClient computev1connect.ComputeServiceClient
}

func NewFinOpsServer(computeURL string) *FinOpsServer {
	return &FinOpsServer{
		computeClient: computev1connect.NewComputeServiceClient(http.DefaultClient, computeURL),
	}
}

func (s *FinOpsServer) ValidateBudget(ctx context.Context, req *connect.Request[finopsv1.ValidateBudgetRequest]) (*connect.Response[finopsv1.ValidateBudgetResponse], error) {
	slog.Info("ValidateBudget", "project", req.Msg.ProjectId, "amount", req.Msg.RequestedAmount)
	
	// Deep Emulation: Query Compute for service health before approving budget
	healthRes, err := s.computeClient.CheckHealth(ctx, connect.NewRequest(&computev1.CheckHealthRequest{
		ServiceName: "finops-validator",
	}))
	
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to reach compute service: %w", err))
	}

	if healthRes.Msg.Status != computev1.CheckHealthResponse_HEALTHY {
		return connect.NewResponse(&finopsv1.ValidateBudgetResponse{
			Approved: false,
			Message:  "Compute substrate unhealthy, budget approval suspended",
		}), nil
	}

	return connect.NewResponse(&finopsv1.ValidateBudgetResponse{
		Approved: true,
		Message:  "Budget approved against healthy compute substrate",
	}), nil
}

func (s *FinOpsServer) EstimateCost(ctx context.Context, req *connect.Request[finopsv1.EstimateCostRequest]) (*connect.Response[finopsv1.EstimateCostResponse], error) {
	slog.Info("EstimateCost", "service", req.Msg.Service, "action", req.Msg.Action)
	return connect.NewResponse(&finopsv1.EstimateCostResponse{
		EstimatedUsd:  0.05,
		Confidence:    "LOW",
	}), nil
}

func (s *FinOpsServer) TrackUsage(ctx context.Context, req *connect.Request[finopsv1.TrackUsageRequest]) (*connect.Response[finopsv1.TrackUsageResponse], error) {
	slog.Info("TrackUsage", "service", req.Msg.Service, "units", req.Msg.ConsumptionUnits)
	return connect.NewResponse(&finopsv1.TrackUsageResponse{CurrentMtdUsd: 0.0}), nil
}
