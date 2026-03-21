package inference

import (
	"context"
	"testing"

	computev1 "olympus.fleet/00SDLC/OlympusGCP-Compute/40000-Communication-Contracts/40400-Protocol-Synthetics/connect-rpc/gen/v1/compute"
	"connectrpc.com/connect"
)

func TestComputeServer_CoverageExpansion(t *testing.T) {
	server := &ComputeServer{}
	ctx := context.Background()

	// 1. Test RunService
	res, err := server.RunService(ctx, connect.NewRequest(&computev1.RunServiceRequest{
		ServiceName: "worker-1",
		Image: "gcr.io/test/worker:latest",
	}))
	if err != nil {
		t.Fatalf("RunService failed: %v", err)
	}
	if res.Msg.EndpointUrl == "" {
		t.Error("Expected endpoint URL")
	}

	// 2. Test TriggerFunction
	funcRes, err := server.TriggerFunction(ctx, connect.NewRequest(&computev1.TriggerFunctionRequest{
		FunctionName: "handler",
		Payload: `{"data":"test"}`,
	}))
	if err != nil {
		t.Fatalf("TriggerFunction failed: %v", err)
	}
	if funcRes.Msg.Result == "" {
		t.Error("Expected result")
	}
}
