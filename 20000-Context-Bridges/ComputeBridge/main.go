package main

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"github.com/mark3labs/mcp-go/mcp"

	"olympus.fleet/00SDLC/OlympusGCP-Compute/40000-Communication-Contracts/40400-Protocol-Synthetics/connect-rpc/gen/v1/compute/computev1connect"
	computev1 "olympus.fleet/00SDLC/OlympusGCP-Compute/40000-Communication-Contracts/40400-Protocol-Synthetics/connect-rpc/gen/v1/compute"
	"olympus.fleet/00SDLC/Olympus2/90000-Enablement-Labs/90200-Logic-Libraries/140-MCPBridge"
)

func main() {
	s := mcpbridge.NewBridgeServer("OlympusComputeBridge", "1.0.0")

	client := computev1connect.NewComputeServiceClient(
		http.DefaultClient,
		"http://localhost:8095",
	)

	s.AddTool(mcp.NewTool("compute_run_service",
		mcp.WithDescription("Launch a compute service. Args: {name: string, image: string}"),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		m, err := mcpbridge.ExtractMap(request)
		if err != nil {
			return mcpbridge.HandleError(err)
		}

		name, _ := m["name"].(string)
		image, _ := m["image"].(string)

		resp, err := client.RunService(ctx, connect.NewRequest(&computev1.RunServiceRequest{
			ServiceName: name,
			Image:       image,
		}))
		if err != nil {
			return mcpbridge.HandleError(err)
		}

		return mcp.NewToolResultText(fmt.Sprintf("Service '%s' running at: %s", name, resp.Msg.EndpointUrl)), nil
	})

	s.Run()
}
