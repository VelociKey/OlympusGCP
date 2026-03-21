package main

import (
	"context"
	"testing"

	eventsv1 "olympus.fleet/00SDLC/OlympusGCP-Events/40000-Communication-Contracts/40400-Protocol-Synthetics/connect-rpc/gen/v1/events"
	"olympus.fleet/00SDLC/OlympusGCP-Events/10000-Autonomous-Actors/10700-Processing-Engines/10710-Reasoning-Inference/inference"
	"connectrpc.com/connect"
)

func TestEventsServerAdvanced(t *testing.T) {
	server := inference.NewEventsServer()
	ctx := context.Background()

	// Test Pause/Resume
	queue := "test-queue"
	server.PauseQueue(ctx, connect.NewRequest(&eventsv1.PauseQueueRequest{Name: queue}))
	
	_, err := server.CreateTask(ctx, connect.NewRequest(&eventsv1.CreateTaskRequest{
		Queue: queue,
		Payload: "blocked",
	}))
	if err == nil {
		t.Error("Expected error for paused queue, got nil")
	}

	server.ResumeQueue(ctx, connect.NewRequest(&eventsv1.ResumeQueueRequest{Name: queue}))
	_, err = server.CreateTask(ctx, connect.NewRequest(&eventsv1.CreateTaskRequest{
		Queue: queue,
		Payload: "allowed",
	}))
	if err != nil {
		t.Fatalf("CreateTask failed after resume: %v", err)
	}

	// Test Deduplication
	taskID := "unique-123"
	server.CreateTask(ctx, connect.NewRequest(&eventsv1.CreateTaskRequest{
		Queue: queue,
		TaskId: taskID,
		Payload: "first",
	}))
	
	res, _ := server.CreateTask(ctx, connect.NewRequest(&eventsv1.CreateTaskRequest{
		Queue: queue,
		TaskId: taskID,
		Payload: "second",
	}))
	if res.Msg.TaskName != taskID {
		t.Errorf("Expected deduplicated task name %s, got %s", taskID, res.Msg.TaskName)
	}

	// Test Purge
	server.PurgeQueue(ctx, connect.NewRequest(&eventsv1.PurgeQueueRequest{Name: queue}))
	// After purge, same task ID should be allowed again (simulated logic)
	res2, _ := server.CreateTask(ctx, connect.NewRequest(&eventsv1.CreateTaskRequest{
		Queue: queue,
		TaskId: taskID,
		Payload: "third",
	}))
	if res2.Msg.TaskName != taskID {
		t.Errorf("Expected new task after purge, got %s", res2.Msg.TaskName)
	}
}

func TestEventsServerBasic(t *testing.T) {
	server := inference.NewEventsServer()
	ctx := context.Background()

	// Test Publish
	pubRes, err := server.Publish(ctx, connect.NewRequest(&eventsv1.PublishRequest{
		Topic: "test-topic",
		Data:  "test-data",
	}))
	if err != nil {
		t.Fatalf("Publish failed: %v", err)
	}
	if pubRes.Msg.MessageId == "" {
		t.Error("Expected message ID, got empty string")
	}

	// Test CreateJob
	jobRes, err := server.CreateJob(ctx, connect.NewRequest(&eventsv1.CreateJobRequest{
		Name:     "test-job",
		Schedule: "* * * * *",
	}))
	if err != nil {
		t.Fatalf("CreateJob failed: %v", err)
	}
	if jobRes.Msg.JobId == "" {
		t.Error("Expected job ID, got empty string")
	}
}
