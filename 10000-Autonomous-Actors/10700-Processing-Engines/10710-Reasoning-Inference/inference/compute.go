package inference

import (
	"context"
	"fmt"
	"log/slog"

	computev1 "olympus.fleet/00SDLC/OlympusGCP-Compute/40000-Communication-Contracts/40400-Protocol-Synthetics/connect-rpc/gen/v1/compute"
	"connectrpc.com/connect"
)

type ComputeServer struct{}

func (s *ComputeServer) RunService(ctx context.Context, req *connect.Request[computev1.RunServiceRequest]) (*connect.Response[computev1.RunServiceResponse], error) {
	slog.Info("RunService", "name", req.Msg.ServiceName, "image", req.Msg.Image)
	endpoint := fmt.Sprintf("http://localhost:8080/services/%s", req.Msg.ServiceName)
	return connect.NewResponse(&computev1.RunServiceResponse{EndpointUrl: endpoint}), nil
}

func (s *ComputeServer) TriggerFunction(ctx context.Context, req *connect.Request[computev1.TriggerFunctionRequest]) (*connect.Response[computev1.TriggerFunctionResponse], error) {
	slog.Info("TriggerFunction", "name", req.Msg.FunctionName)
	result := fmt.Sprintf("Function %s executed successfully.", req.Msg.FunctionName)
	return connect.NewResponse(&computev1.TriggerFunctionResponse{Result: result}), nil
}

func (s *ComputeServer) CheckHealth(ctx context.Context, req *connect.Request[computev1.CheckHealthRequest]) (*connect.Response[computev1.CheckHealthResponse], error) {
	slog.Info("CheckHealth", "name", req.Msg.ServiceName)
	return connect.NewResponse(&computev1.CheckHealthResponse{
		Status:  computev1.CheckHealthResponse_HEALTHY,
		Message: "Service is operational",
	}), nil
}
