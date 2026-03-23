package inference

import (
	"context"
	"testing"

	messagingv1 "olympus.fleet/00SDLC/OlympusGCP/gen/messaging"
	"connectrpc.com/connect"
)

func TestMessagingServer_CoverageExpansion(t *testing.T) {
	server := &MessagingServer{}
	ctx := context.Background()

	// 1. Test SendMessage
	res, err := server.SendMessage(ctx, connect.NewRequest(&messagingv1.SendMessageRequest{
		Recipient: "user1",
		Content: "msg",
	}))
	if err != nil || res.Msg.MessageId == "" {
		t.Error("SendMessage failed")
	}
}
