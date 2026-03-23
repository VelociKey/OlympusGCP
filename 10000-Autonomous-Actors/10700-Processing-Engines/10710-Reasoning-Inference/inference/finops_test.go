package inference

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	finopsv1 "olympus.fleet/00SDLC/OlympusGCP/gen/finops"
	computev1 "olympus.fleet/00SDLC/OlympusGCP/gen/compute"
	"olympus.fleet/00SDLC/OlympusGCP/gen/google/cloud/compute/computev1connect"
	"connectrpc.com/connect"
)

type mockComputeHandler struct {
	computev1connect.UnimplementedComputeServiceHandler
}

func (h *mockComputeHandler) CheckHealth(ctx context.Context, req *connect.Request[computev1.CheckHealthRequest]) (*connect.Response[computev1.CheckHealthResponse], error) {
	return connect.NewResponse(&computev1.CheckHealthResponse{Status: computev1.CheckHealthResponse_HEALTHY}), nil
}

func TestFinOpsServer_CoverageExpansion(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle(computev1connect.NewComputeServiceHandler(&mockComputeHandler{}))
	srv := httptest.NewServer(mux)
	defer srv.Close()

	server := NewFinOpsServer(srv.URL)
	ctx := context.Background()

	// 1. Test ValidateBudget
	res, err := server.ValidateBudget(ctx, connect.NewRequest(&finopsv1.ValidateBudgetRequest{
		ProjectId: "p1",
		RequestedAmount: 50.0,
	}))
	if err != nil || !res.Msg.Approved {
		t.Errorf("ValidateBudget failed: %v", err)
	}

	// 2. Test EstimateCost
	_, err = server.EstimateCost(ctx, connect.NewRequest(&finopsv1.EstimateCostRequest{
		Service: "storage",
		Action: "run",
	}))
	if err != nil {
		t.Error("EstimateCost failed")
	}

	// 3. Test TrackUsage
	_, err = server.TrackUsage(ctx, connect.NewRequest(&finopsv1.TrackUsageRequest{
		Service: "storage",
		ConsumptionUnits: 5,
	}))
	if err != nil {
		t.Error("TrackUsage failed")
	}
}
