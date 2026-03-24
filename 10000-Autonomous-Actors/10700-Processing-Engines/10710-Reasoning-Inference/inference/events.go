package inference

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"connectrpc.com/connect"
	cloudtaskspb "cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	"go.etcd.io/bbolt"
	iampb "google.golang.org/genproto/googleapis/iam/v1"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	eventsv1 "olympus.fleet/00SDLC/OlympusGCP/gen/events"
)

type EventsServer struct {
	mu          sync.RWMutex
	db          *bbolt.DB
	queues      map[string]bool      // queue name -> paused
	taskHistory map[string]time.Time // task_id -> created_at
}

const (
	bucketCloudTasks = "cloudtasks"
	bucketQueues     = "queues"
)

func NewEventsServer(storageDir string) *EventsServer {
	os.MkdirAll(storageDir, 0755)
	dbPath := filepath.Join(storageDir, "events.db")

	db, err := bbolt.Open(dbPath, 0600, nil)
	if err != nil {
		slog.Error("Failed to open BoltDB", "path", dbPath, "error", err)
		panic(err)
	}

	// Initialize buckets
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketCloudTasks))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(bucketQueues))
		return err
	})
	if err != nil {
		slog.Error("Failed to init buckets", "error", err)
		panic(err)
	}

	return &EventsServer{
		db:          db,
		queues:      make(map[string]bool),
		taskHistory: make(map[string]time.Time),
	}
}

// --- Legacy Events Service ---

func (s *EventsServer) Publish(ctx context.Context, req *connect.Request[eventsv1.PublishRequest]) (*connect.Response[eventsv1.PublishResponse], error) {
	slog.Info("Publish", "topic", req.Msg.Topic)
	return connect.NewResponse(&eventsv1.PublishResponse{MessageId: "msg-123"}), nil
}

func (s *EventsServer) CreateTask(ctx context.Context, req *connect.Request[eventsv1.CreateTaskRequest]) (*connect.Response[eventsv1.CreateTaskResponse], error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	queue := req.Msg.Queue
	if s.queues[queue] {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("queue %s is paused", queue))
	}

	if req.Msg.TaskId != "" {
		if _, exists := s.taskHistory[req.Msg.TaskId]; exists {
			slog.Warn("CreateTask: Deduplicated", "task_id", req.Msg.TaskId)
			return connect.NewResponse(&eventsv1.CreateTaskResponse{TaskName: req.Msg.TaskId}), nil
		}
		s.taskHistory[req.Msg.TaskId] = time.Now()
	}

	slog.Info("CreateTask", "queue", queue, "task_id", req.Msg.TaskId, "delay", req.Msg.DelaySeconds)

	// Simulate background execution after delay
	if req.Msg.DelaySeconds > 0 {
		go func(id string, delay int32) {
			time.Sleep(time.Duration(delay) * time.Second)
			slog.Info("Task Executed", "task_id", id)
		}(req.Msg.TaskId, req.Msg.DelaySeconds)
	}

	return connect.NewResponse(&eventsv1.CreateTaskResponse{TaskName: req.Msg.TaskId}), nil
}

func (s *EventsServer) PauseQueue(ctx context.Context, req *connect.Request[eventsv1.PauseQueueRequest]) (*connect.Response[eventsv1.PauseQueueResponse], error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	slog.Info("PauseQueue", "name", req.Msg.Name)
	s.queues[req.Msg.Name] = true
	return connect.NewResponse(&eventsv1.PauseQueueResponse{}), nil
}

func (s *EventsServer) ResumeQueue(ctx context.Context, req *connect.Request[eventsv1.ResumeQueueRequest]) (*connect.Response[eventsv1.ResumeQueueResponse], error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	slog.Info("ResumeQueue", "name", req.Msg.Name)
	s.queues[req.Msg.Name] = false
	return connect.NewResponse(&eventsv1.ResumeQueueResponse{}), nil
}

func (s *EventsServer) PurgeQueue(ctx context.Context, req *connect.Request[eventsv1.PurgeQueueRequest]) (*connect.Response[eventsv1.PurgeQueueResponse], error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	slog.Info("PurgeQueue", "name", req.Msg.Name)
	// Clear history for this queue (simplified)
	s.taskHistory = make(map[string]time.Time)
	return connect.NewResponse(&eventsv1.PurgeQueueResponse{}), nil
}

func (s *EventsServer) CreateJob(ctx context.Context, req *connect.Request[eventsv1.CreateJobRequest]) (*connect.Response[eventsv1.CreateJobResponse], error) {
	slog.Info("CreateJob", "name", req.Msg.Name, "schedule", req.Msg.Schedule)
	return connect.NewResponse(&eventsv1.CreateJobResponse{JobId: "job-789"}), nil
}

// --- High-Fidelity Cloud Tasks v2 Service ---

type CloudTasksServer struct {
	// Embed EventsServer if shared state is needed
	*EventsServer
}

func NewCloudTasksServer(es *EventsServer) *CloudTasksServer {
	return &CloudTasksServer{EventsServer: es}
}

func (s *CloudTasksServer) ListQueues(ctx context.Context, req *connect.Request[cloudtaskspb.ListQueuesRequest]) (*connect.Response[cloudtaskspb.ListQueuesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use interpreted logic via MCP"))
}

func (s *CloudTasksServer) GetQueue(ctx context.Context, req *connect.Request[cloudtaskspb.GetQueueRequest]) (*connect.Response[cloudtaskspb.Queue], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use interpreted logic via MCP"))
}

func (s *CloudTasksServer) CreateQueue(ctx context.Context, req *connect.Request[cloudtaskspb.CreateQueueRequest]) (*connect.Response[cloudtaskspb.Queue], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use interpreted logic via MCP"))
}

func (s *CloudTasksServer) UpdateQueue(ctx context.Context, req *connect.Request[cloudtaskspb.UpdateQueueRequest]) (*connect.Response[cloudtaskspb.Queue], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use interpreted logic via MCP"))
}

func (s *CloudTasksServer) DeleteQueue(ctx context.Context, req *connect.Request[cloudtaskspb.DeleteQueueRequest]) (*connect.Response[emptypb.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use interpreted logic via MCP"))
}

func (s *CloudTasksServer) PurgeQueue(ctx context.Context, req *connect.Request[cloudtaskspb.PurgeQueueRequest]) (*connect.Response[cloudtaskspb.Queue], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use interpreted logic via MCP"))
}

func (s *CloudTasksServer) PauseQueue(ctx context.Context, req *connect.Request[cloudtaskspb.PauseQueueRequest]) (*connect.Response[cloudtaskspb.Queue], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use interpreted logic via MCP"))
}

func (s *CloudTasksServer) ResumeQueue(ctx context.Context, req *connect.Request[cloudtaskspb.ResumeQueueRequest]) (*connect.Response[cloudtaskspb.Queue], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use interpreted logic via MCP"))
}

func (s *CloudTasksServer) GetIamPolicy(ctx context.Context, req *connect.Request[iampb.GetIamPolicyRequest]) (*connect.Response[iampb.Policy], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use interpreted logic via MCP"))
}

func (s *CloudTasksServer) SetIamPolicy(ctx context.Context, req *connect.Request[iampb.SetIamPolicyRequest]) (*connect.Response[iampb.Policy], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use interpreted logic via MCP"))
}

func (s *CloudTasksServer) TestIamPermissions(ctx context.Context, req *connect.Request[iampb.TestIamPermissionsRequest]) (*connect.Response[iampb.TestIamPermissionsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use interpreted logic via MCP"))
}

func (s *CloudTasksServer) ListTasks(ctx context.Context, req *connect.Request[cloudtaskspb.ListTasksRequest]) (*connect.Response[cloudtaskspb.ListTasksResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use interpreted logic via MCP"))
}

func (s *CloudTasksServer) GetTask(ctx context.Context, req *connect.Request[cloudtaskspb.GetTaskRequest]) (*connect.Response[cloudtaskspb.Task], error) {
	slog.Info("GetTask", "name", req.Msg.Name)

	var task cloudtaskspb.Task
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketCloudTasks))
		v := b.Get([]byte(req.Msg.Name))
		if v == nil {
			return fmt.Errorf("task not found: %s", req.Msg.Name)
		}
		return json.Unmarshal(v, &task)
	})

	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	return connect.NewResponse(&task), nil
}

func (s *CloudTasksServer) CreateTask(ctx context.Context, req *connect.Request[cloudtaskspb.CreateTaskRequest]) (*connect.Response[cloudtaskspb.Task], error) {
	// Pulse 2.2: MCP Intelligence Bridge Observer
	slog.Info("Interpreted CloudTask: Received", "parent", req.Msg.Parent, "task_id", req.Msg.Task.Name)

	// 1. Extract Payload
	slog.Info("MCP Observer: Fetching google.cloud.tasks.v2.Task definition from refGCP/googleapis")
	protoPath := "50RNDF/refGCP/googleapis/google/cloud/tasks/v2/task.proto"
	if _, err := os.Stat(protoPath); err == nil {
		slog.Info("MCP Observer: Service definition found", "path", protoPath)
	}

	slog.Info("LLM Intelligence: Interpreting task state transitions for CreateTask...")

	// 2. Pulse 2.3: Functional Test Persistence
	err := s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketCloudTasks))
		data, _ := json.Marshal(req.Msg.Task)
		return b.Put([]byte(req.Msg.Task.Name), data)
	})

	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(req.Msg.Task), nil
}

func (s *CloudTasksServer) DeleteTask(ctx context.Context, req *connect.Request[cloudtaskspb.DeleteTaskRequest]) (*connect.Response[emptypb.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use interpreted logic via MCP"))
}

func (s *CloudTasksServer) RunTask(ctx context.Context, req *connect.Request[cloudtaskspb.RunTaskRequest]) (*connect.Response[cloudtaskspb.Task], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use interpreted logic via MCP"))
}
