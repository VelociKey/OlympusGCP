package main

import (
	"context"
	"testing"

	messagingv1 "olympus.fleet/00SDLC/OlympusGCP-Messaging/40000-Communication-Contracts/40400-Protocol-Synthetics/connect-rpc/gen/v1/messaging"
	"olympus.fleet/00SDLC/OlympusGCP-Messaging/10000-Autonomous-Actors/10700-Processing-Engines/10710-Reasoning-Inference/inference"
	"connectrpc.com/connect"
)

func TestMessagingServer(t *testing.T) {
	server := &inference.MessagingServer{}
	ctx := context.Background()

	// Test SendMessage
	req := connect.NewRequest(&messagingv1.SendMessageRequest{
		Recipient: "user-1",
		Content:   "Hello",
	})
	res, err := server.SendMessage(ctx, req)
	if err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}
	if res.Msg.MessageId == "" {
		t.Error("Expected message ID, got empty string")
	}
}
