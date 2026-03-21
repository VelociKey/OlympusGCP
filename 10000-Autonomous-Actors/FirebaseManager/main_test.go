package main

import (
	"context"
	"testing"

	firebasev1 "olympus.fleet/00SDLC/OlympusGCP-Firebase/40000-Communication-Contracts/40400-Protocol-Synthetics/connect-rpc/gen/v1/firebase"
	"olympus.fleet/00SDLC/OlympusGCP-Firebase/10000-Autonomous-Actors/10700-Processing-Engines/10710-Reasoning-Inference/inference"
	"connectrpc.com/connect"
)

func TestFirebaseServer(t *testing.T) {
	server := &inference.FirebaseServer{}
	ctx := context.Background()

	// Test CreateUser
	userRes, err := server.CreateUser(ctx, connect.NewRequest(&firebasev1.CreateUserRequest{
		Email: "test@example.com",
	}))
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	if userRes.Msg.Uid == "" {
		t.Error("Expected UID, got empty string")
	}

	// Test SetDocument
	docRes, err := server.SetDocument(ctx, connect.NewRequest(&firebasev1.SetDocumentRequest{
		Collection: "data",
		DocId:      "123",
		DataJson:   `{"foo": "bar"}`,
	}))
	if err != nil {
		t.Fatalf("SetDocument failed: %v", err)
	}
	if docRes.Msg.Uid != "123" {
		t.Errorf("Expected DocID '123', got '%s'", docRes.Msg.Uid)
	}
}
