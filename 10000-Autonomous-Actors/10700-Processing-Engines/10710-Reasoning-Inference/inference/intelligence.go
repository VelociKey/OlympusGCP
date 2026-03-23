package inference

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	intelligencev1 "olympus.fleet/00SDLC/OlympusGCP/gen/intelligence"
	"connectrpc.com/connect"
)

type IntelligenceServer struct{}

func (s *IntelligenceServer) Predict(ctx context.Context, req *connect.Request[intelligencev1.PredictRequest]) (*connect.Response[intelligencev1.PredictResponse], error) {
	prompt := req.Msg.Prompt
	slog.Info("Predict", "model", req.Msg.Model, "prompt_len", len(prompt))

	// Deep Emulation: Validate Prompt Format
	// Requirement: Prompt must start with 'MISSION:' or 'QUERY:' to be considered valid for Sovereign processing.
	valid := false
	if strings.HasPrefix(strings.ToUpper(prompt), "MISSION:") || strings.HasPrefix(strings.ToUpper(prompt), "QUERY:") {
		valid = true
	}

	if !valid {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid prompt grammar: must start with MISSION: or QUERY:"))
	}

	prediction := fmt.Sprintf("[Sovereign Intelligence] Model: %s processed valid request. Output: Analyzed.", req.Msg.Model)
	return connect.NewResponse(&intelligencev1.PredictResponse{Prediction: prediction}), nil
}
