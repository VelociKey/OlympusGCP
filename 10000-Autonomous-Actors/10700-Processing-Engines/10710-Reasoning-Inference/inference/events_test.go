package inference

import (
	"context"
	"testing"

	eventsv1 "olympus.fleet/00SDLC/OlympusGCP/gen/events"
	"connectrpc.com/connect"
)

func TestEventsServer_CoverageExpansion(t *testing.T) {
	server := NewEventsServer()
	ctx := context.Background()

	// 1. Test CreateTask with Background Execution
	res, err := server.CreateTask(ctx, connect.NewRequest(&eventsv1.CreateTaskRequest{
		Queue: "default",
		TaskId: "task-bg",
		DelaySeconds: 1,
	}))
	if err != nil {
		t.Fatalf("CreateTask failed: %v", err)
	}
	if res.Msg.TaskName != "task-bg" {
		t.Errorf("Expected task name 'task-bg', got %s", res.Msg.TaskName)
	}

	// 2. Test Resume Queue (which is already unpaused by default)
	server.ResumeQueue(ctx, connect.NewRequest(&eventsv1.ResumeQueueRequest{Name: "default"}))
	
	// 3. Test Purge Queue
	server.PurgeQueue(ctx, connect.NewRequest(&eventsv1.PurgeQueueRequest{Name: "default"}))

	// 4. Test Publish
	pubRes, err := server.Publish(ctx, connect.NewRequest(&eventsv1.PublishRequest{
		Topic: "test",
		Data: "data",
	}))
	if err != nil {
		t.Fatalf("Publish failed: %v", err)
	}
	if pubRes.Msg.MessageId == "" {
		t.Error("Expected message ID")
	}
}
