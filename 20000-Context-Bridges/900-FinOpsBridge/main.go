package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/protobuf/reflect/protoreflect"

	mcpv1 "olympus.fleet/00SDLC/Olympus2/40000-Communication-Contracts/40400-Protocol-Synthetics/connect-rpc/olympus/mcp/v1"
	mcpv1connect "olympus.fleet/00SDLC/Olympus2/40000-Communication-Contracts/40400-Protocol-Synthetics/connect-rpc/olympus/mcp/v1/mcpv1connect"

	finopsv1 "olympus.fleet/00SDLC/OlympusGCP-FinOps/40000-Communication-Contracts/40400-Protocol-Synthetics/connect-rpc/gen/v1/finops"
	"olympus.fleet/00SDLC/OlympusGCP-FinOps/40000-Communication-Contracts/40400-Protocol-Synthetics/connect-rpc/gen/v1/finops/finopsv1connect"
)

type FinOpsBridgeServer struct {
	client finopsv1connect.FinOpsServiceClient
	logger *slog.Logger
}

// ---------------------------------------------------------
// Tools Implementation
// ---------------------------------------------------------
func (s *FinOpsBridgeServer) ListTools(
	ctx context.Context,
	req *connect.Request[mcpv1.ListToolsRequest],
) (*connect.Response[mcpv1.ListToolsResponse], error) {

	estCostSchema, _ := structpb.NewStruct(map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"service": map[string]interface{}{"type": "string"},
			"action":  map[string]interface{}{"type": "string"},
		},
		"required": []interface{}{"service", "action"},
	})

	valBudgetSchema, _ := structpb.NewStruct(map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"amount": map[string]interface{}{"type": "number"},
		},
		"required": []interface{}{"amount"},
	})

	// Define the tools this bridge exposes
	tools := []*mcpv1.Tool{
		{
			Name:        "finops_estimate_cost",
			Description: "Estimate the GCP cost of a request based on the service and action.",
			InputSchema: estCostSchema,
		},
		{
			Name:        "finops_validate_budget",
			Description: "Check if an operation fits within the local safety budget limit.",
			InputSchema: valBudgetSchema,
		},
	}

	return connect.NewResponse(&mcpv1.ListToolsResponse{
		Tools: tools,
	}), nil
}

func (s *FinOpsBridgeServer) CallTool(
	ctx context.Context,
	req *connect.Request[mcpv1.CallToolRequest],
) (*connect.Response[mcpv1.CallToolResponse], error) {

	args := req.Msg.Arguments.AsMap()

	switch req.Msg.Name {
	case "finops_estimate_cost":
		service, _ := args["service"].(string)
		action, _ := args["action"].(string)

		resp, err := s.client.EstimateCost(ctx, connect.NewRequest(&finopsv1.EstimateCostRequest{
			Service: service,
			Action:  action,
		}))

		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		// Return structurally typed jeBNF payload
		msg := fmt.Sprintf("Result { EstimatedCost = %.2f; Confidence = \"%s\"; }", resp.Msg.EstimatedUsd, resp.Msg.Confidence)
		return connect.NewResponse(&mcpv1.CallToolResponse{
			Content: []*mcpv1.Content{{Type: "olympus.fleet/00SDLC/Olympus2/01000-Identity-Foundations/P0000-pkg/text", Text: msg}},
		}), nil

	case "finops_validate_budget":
		amount, _ := args["amount"].(float64)

		resp, err := s.client.ValidateBudget(ctx, connect.NewRequest(&finopsv1.ValidateBudgetRequest{
			ProjectId:       "local-dev",
			RequestedAmount: amount,
		}))

		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		// Return structurally typed jeBNF payload
		msg := fmt.Sprintf("Result { Approved = %t; Message = \"%s\"; }", resp.Msg.Approved, resp.Msg.Message)
		return connect.NewResponse(&mcpv1.CallToolResponse{
			Content: []*mcpv1.Content{{Type: "olympus.fleet/00SDLC/Olympus2/01000-Identity-Foundations/P0000-pkg/text", Text: msg}},
		}), nil

	default:
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("tool %s not found", req.Msg.Name))
	}
}

// ---------------------------------------------------------
// required ModelContextProtocolHandler boilerplate (stubbed)
// ---------------------------------------------------------
func (s *FinOpsBridgeServer) Initialize(ctx context.Context, req *connect.Request[mcpv1.InitializeRequest]) (*connect.Response[mcpv1.InitializeResponse], error) {
	return connect.NewResponse(&mcpv1.InitializeResponse{
		ProtocolVersion: "2024-11-05",
		ServerInfo:      &mcpv1.ServerInfo{Name: "FinOpsBridge", Version: "1.0.0"},
		Capabilities:    &mcpv1.ServerCapabilities{Tools: true},
	}), nil
}
func (s *FinOpsBridgeServer) ListResources(context.Context, *connect.Request[mcpv1.ListResourcesRequest]) (*connect.Response[mcpv1.ListResourcesResponse], error) {
	return connect.NewResponse(&mcpv1.ListResourcesResponse{}), nil
}
func (s *FinOpsBridgeServer) ReadResource(context.Context, *connect.Request[mcpv1.ReadResourceRequest]) (*connect.Response[mcpv1.ReadResourceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("unimplemented"))
}
func (s *FinOpsBridgeServer) ListResourceTemplates(context.Context, *connect.Request[mcpv1.ListResourceTemplatesRequest]) (*connect.Response[mcpv1.ListResourceTemplatesResponse], error) {
	return connect.NewResponse(&mcpv1.ListResourceTemplatesResponse{}), nil
}
func (s *FinOpsBridgeServer) ListPrompts(context.Context, *connect.Request[mcpv1.ListPromptsRequest]) (*connect.Response[mcpv1.ListPromptsResponse], error) {
	return connect.NewResponse(&mcpv1.ListPromptsResponse{}), nil
}
func (s *FinOpsBridgeServer) GetPrompt(context.Context, *connect.Request[mcpv1.GetPromptRequest]) (*connect.Response[mcpv1.GetPromptResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("unimplemented"))
}

// ---------------------------------------------------------
// MAIN
// ---------------------------------------------------------
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	// Connect to internal Manager
	backend := finopsv1connect.NewFinOpsServiceClient(http.DefaultClient, "http://localhost:8093")

	server := &FinOpsBridgeServer{
		client: backend,
		logger: logger,
	}

	mux := http.NewServeMux()
	path, handler := mcpv1connect.NewModelContextProtocolHandler(server)
	mux.Handle(path, handler)

	// Bind HTTP2/Connect endpoint (e.g. 8094 to avoid colliding with Manager)
	addr := "127.0.0.1:8094"
	logger.Info("Starting FinOps Connect RPC Bridge", "addr", addr)

	// Note: Explicit ReadHeaderTimeout to satisfy G114!
	srv := &http.Server{
		Addr:              addr,
		Handler:           h2c.NewHandler(mux, &http2.Server{}),
		ReadHeaderTimeout: 5 * 1000 * 1000 * 1000,
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Error("bridge stopped", "error", err)
		os.Exit(1)
	}
}
