package inference

import (
	"context"
	"testing"

	intelligencev1 "olympus.fleet/00SDLC/OlympusGCP/gen/intelligence"
	"connectrpc.com/connect"
)

func TestIntelligenceServer_CoverageExpansion(t *testing.T) {
	server := &IntelligenceServer{}
	ctx := context.Background()

	// 1. Test Predict
	res, err := server.Predict(ctx, connect.NewRequest(&intelligencev1.PredictRequest{
		Model: "gemini",
		Prompt: "MISSION: test prompt",
	}))
	if err != nil || res.Msg.Prediction == "" {
		t.Error("Predict failed")
	}
}
