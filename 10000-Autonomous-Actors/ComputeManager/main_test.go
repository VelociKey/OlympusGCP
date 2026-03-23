package main

import (
	"context"
	"testing"

	computev1 "olympus.fleet/00SDLC/OlympusGCP/gen/compute"
	"olympus.fleet/00SDLC/OlympusGCP/10000-Autonomous-Actors/10700-Processing-Engines/10710-Reasoning-Inference/inference"
	"connectrpc.com/connect"
)

func TestComputeServerAdvanced(t *testing.T) {
	server := &inference.ComputeServer{}
	ctx := context.Background()

	// Test CheckHealth
	req := connect.NewRequest(&computev1.CheckHealthRequest{ServiceName: "api"})
	res, err := server.CheckHealth(ctx, req)
	if err != nil {
		t.Fatalf("CheckHealth failed: %v", err)
	}
	if res.Msg.Status != computev1.CheckHealthResponse_HEALTHY {
		t.Errorf("Expected HEALTHY status, got %v", res.Msg.Status)
	}
}
