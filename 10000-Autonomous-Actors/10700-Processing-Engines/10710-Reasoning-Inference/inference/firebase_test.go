package inference

import (
	"context"
	"testing"

	firebasev1 "olympus.fleet/00SDLC/OlympusGCP/gen/firebase"
	"connectrpc.com/connect"
)

func TestFirebaseServer_CoverageExpansion(t *testing.T) {
	server := &FirebaseServer{}
	ctx := context.Background()

	// 1. Test CreateUser
	res, err := server.CreateUser(ctx, connect.NewRequest(&firebasev1.CreateUserRequest{
		Email: "test@example.com",
	}))
	if err != nil || res.Msg.Uid == "" {
		t.Error("CreateUser failed")
	}

	// 2. Test SetDocument
	_, err = server.SetDocument(ctx, connect.NewRequest(&firebasev1.SetDocumentRequest{
		Collection: "users",
		DocId: "u1",
		DataJson: `{}`,
	}))
	if err != nil {
		t.Error("SetDocument failed")
	}
}
