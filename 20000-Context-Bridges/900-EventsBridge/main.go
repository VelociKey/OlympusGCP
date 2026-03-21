package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"

	"connectrpc.com/connect"

	mcpbridge "olympus.fleet/00SDLC/Olympus2/90000-Enablement-Labs/90200-Logic-Libraries/140-MCPBridge"
	eventsv1 "olympus.fleet/00SDLC/OlympusGCP-Events/40000-Communication-Contracts/40400-Protocol-Synthetics/connect-rpc/gen/v1/events"
	"olympus.fleet/00SDLC/OlympusGCP-Events/40000-Communication-Contracts/40400-Protocol-Synthetics/connect-rpc/gen/v1/events/eventsv1connect"
)

func main() {
	s := mcpbridge.NewBridgeServer("OlympusEventsBridge", "1.0.0")

	client := eventsv1connect.NewEventsServiceClient(
		http.DefaultClient,
		"http://localhost:8094",
	)

	s.AddTool(mcp.NewTool("events_publish",
		mcp.WithDescription("Publish a message to a topic. Args: {topic: string, payload: string}"),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		m, err := mcpbridge.ExtractMap(request)
		if err != nil {
			return mcpbridge.HandleError(err)
		}

		topic, _ := m["topic"].(string)
		payload, _ := m["payload"].(string)

		resp, err := client.Publish(ctx, connect.NewRequest(&eventsv1.PublishRequest{
			Topic: topic,
			Data:  payload,
		}))
		if err != nil {
			return mcpbridge.HandleError(err)
		}

		return mcp.NewToolResultText(fmt.Sprintf("Message published. ID: %s", resp.Msg.MessageId)), nil
	})

	s.AddTool(mcp.NewTool("events_create_task",
		mcp.WithDescription("Create a background task. Args: {queue: string, payload: string}"),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		m, err := mcpbridge.ExtractMap(request)
		if err != nil {
			return mcpbridge.HandleError(err)
		}

		queue, _ := m["queue"].(string)
		payload, _ := m["payload"].(string)

		resp, err := client.CreateTask(ctx, connect.NewRequest(&eventsv1.CreateTaskRequest{
			Queue:   queue,
			Payload: payload,
		}))
		if err != nil {
			return mcpbridge.HandleError(err)
		}

		return mcp.NewToolResultText(fmt.Sprintf("Task created: %s", resp.Msg.TaskName)), nil
	})

	s.Run()
}
