package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	finopsv1 "olympus.fleet/00SDLC/OlympusGCP-FinOps/40000-Communication-Contracts/40400-Protocol-Synthetics/connect-rpc/gen/v1/finops"
	computev1 "olympus.fleet/00SDLC/OlympusGCP-Compute/40000-Communication-Contracts/40400-Protocol-Synthetics/connect-rpc/gen/v1/compute"
	"olympus.fleet/00SDLC/OlympusGCP-Compute/40000-Communication-Contracts/40400-Protocol-Synthetics/connect-rpc/gen/v1/compute/computev1connect"
	"olympus.fleet/00SDLC/OlympusGCP-FinOps/10000-Autonomous-Actors/10700-Processing-Engines/10710-Reasoning-Inference/inference"
	"connectrpc.com/connect"
)

type mockComputeHandler struct {
	computev1connect.UnimplementedComputeServiceHandler
}

func (h *mockComputeHandler) CheckHealth(ctx context.Context, req *connect.Request[computev1.CheckHealthRequest]) (*connect.Response[computev1.CheckHealthResponse], error) {
	return connect.NewResponse(&computev1.CheckHealthResponse{Status: computev1.CheckHealthResponse_HEALTHY}), nil
}

func TestFinOpsServer_DeepEmulation(t *testing.T) {
	// 1. Setup Mock Compute Service
	mux := http.NewServeMux()
	mux.Handle(computev1connect.NewComputeServiceHandler(&mockComputeHandler{}))
	srv := httptest.NewServer(mux)
	defer srv.Close()

	// 2. Initialize FinOps with Mock URL
	server := inference.NewFinOpsServer(srv.URL)
	ctx := context.Background()

	// 3. Test ValidateBudget
	budgetRes, err := server.ValidateBudget(ctx, connect.NewRequest(&finopsv1.ValidateBudgetRequest{
		ProjectId:       "p1",
		RequestedAmount: 100.0,
	}))
	if err != nil {
		t.Fatalf("ValidateBudget failed: %v", err)
	}
	if !budgetRes.Msg.Approved {
		t.Error("Expected budget to be approved against healthy compute")
	}
}
