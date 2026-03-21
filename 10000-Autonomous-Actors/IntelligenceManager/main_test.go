package main

import (
	"context"
	"testing"

	intelligencev1 "olympus.fleet/00SDLC/OlympusGCP-Intelligence/40000-Communication-Contracts/40400-Protocol-Synthetics/connect-rpc/gen/v1/intelligence"
	"olympus.fleet/00SDLC/OlympusGCP-Intelligence/10000-Autonomous-Actors/10700-Processing-Engines/10710-Reasoning-Inference/inference"
	"connectrpc.com/connect"
)

func TestIntelligenceServer(t *testing.T) {
	server := &inference.IntelligenceServer{}
	ctx := context.Background()

	// Test Predict (Valid)
	req := connect.NewRequest(&intelligencev1.PredictRequest{
		Model:  "gemini-pro",
		Prompt: "MISSION: Analyze fleet hardening status",
	})
	res, err := server.Predict(ctx, req)
	if err != nil {
		t.Fatalf("Predict failed: %v", err)
	}
	if res.Msg.Prediction == "" {
		t.Error("Expected prediction, got empty string")
	}

	// Test Predict (Invalid)
	badReq := connect.NewRequest(&intelligencev1.PredictRequest{
		Model:  "gemini-pro",
		Prompt: "Hello world",
	})
	_, err = server.Predict(ctx, badReq)
	if err == nil {
		t.Error("Expected error for invalid prompt format, got nil")
	}
}
