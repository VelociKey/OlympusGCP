package main

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	whisper "olympus.fleet/00SDLC/Olympus2/90000-Enablement-Labs/90200-Logic-Libraries/220-Whisper"
	tasksv1 "olympus.fleet/00SDLC/OlympusGCP/gen/olympus/tasks/v1"
	"olympus.fleet/00SDLC/OlympusGCP/gen/olympus/tasks/v1/tasksv1connect"
)

type TaskServer struct {
	logger *whisper.WhisperLog
	mu     sync.Mutex
	queues map[string][]string
}

func (s *TaskServer) CreateTask(ctx context.Context, req *connect.Request[tasksv1.CreateTaskRequest]) (*connect.Response[tasksv1.CreateTaskResponse], error) {
	start := time.Now()
	taskID := fmt.Sprintf("task-%d", time.Now().UnixNano())

	slog.Info("TaskManager: Creating Task", "queue", req.Msg.QueueId, "target", req.Msg.TargetUrl)

	// Substrate: Local Async Execution
	go func() {
		time.Sleep(2 * time.Second) // Simulate network/queue delay
		slog.Info("TaskManager: Dispatching Task", "id", taskID, "url", req.Msg.TargetUrl)
		resp, err := http.Post(req.Msg.TargetUrl, "application/json", bytes.NewBuffer([]byte(req.Msg.Payload)))
		if err != nil {
			slog.Error("TaskManager: Dispatch Failed", "id", taskID, "error", err)
			return
		}
		defer resp.Body.Close()
		slog.Info("TaskManager: Dispatch Successful", "id", taskID, "status", resp.Status)
	}()

	s.mu.Lock()
	s.queues[req.Msg.QueueId] = append(s.queues[req.Msg.QueueId], taskID)
	s.mu.Unlock()

	s.logger.Log("CREATE_TASK", "SUCCESS", req.Msg.QueueId, taskID, time.Since(start))
	return connect.NewResponse(&tasksv1.CreateTaskResponse{
		TaskId: taskID,
		Status: "ENQUEUED_LOCALLY",
	}), nil
}

func (s *TaskServer) ListQueues(ctx context.Context, req *connect.Request[tasksv1.ListQueuesRequest]) (*connect.Response[tasksv1.ListQueuesResponse], error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var names []string
	for k := range s.queues {
		names = append(names, k)
	}
	return connect.NewResponse(&tasksv1.ListQueuesResponse{QueueNames: names}), nil
}

func main() {
	w := whisper.New("TaskManager", "gcp_tasks.lpsv")
	defer w.Close()

	server := &TaskServer{
		logger: w,
		queues: make(map[string][]string),
	}

	mux := http.NewServeMux()
	path, handler := tasksv1connect.NewTaskServiceHandler(server)
	mux.Handle(path, handler)

	port := "8100"
	slog.Info("TaskManager: Booting Async-Native High-Fidelity Substrate...", "port", port)

	http.ListenAndServe(
		"localhost:"+port,
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
