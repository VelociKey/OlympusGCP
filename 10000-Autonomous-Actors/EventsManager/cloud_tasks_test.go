package main

import (
        "context"
        "net/http"
        "net/http/httptest"
        "testing"

        "connectrpc.com/connect"
        cloudtaskspb "cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
        taskspbconnect "olympus.fleet/00SDLC/OlympusGCP/gen/cloudtasks/apiv2/cloudtaskspb/cloudtaskspbconnect"
        "olympus.fleet/00SDLC/OlympusGCP/10000-Autonomous-Actors/10700-Processing-Engines/10710-Reasoning-Inference/inference"
)

func TestCloudTasksPrototype(t *testing.T) {
        eventsServer := inference.NewEventsServer()
        server := inference.NewCloudTasksServer(eventsServer)
        
        mux := http.NewServeMux()
        path, handler := taskspbconnect.NewCloudTasksHandler(server)
        mux.Handle(path, handler)

        ts := httptest.NewServer(handler)
        defer ts.Close()

        client := taskspbconnect.NewCloudTasksClient(http.DefaultClient, ts.URL)

        t.Run("CreateTask_Interpreted", func(t *testing.T) {
                req := &cloudtaskspb.CreateTaskRequest{
                        Parent: "projects/sovereign-fleet/locations/us-east1/queues/test-queue",
                        Task: &cloudtaskspb.Task{
                                Name: "projects/sovereign-fleet/locations/us-east1/queues/test-queue/tasks/task-123",
                        },
                }

                res, err := client.CreateTask(context.Background(), connect.NewRequest(req))
                if err != nil {
                        t.Fatalf("CreateTask failed: %v", err)
                }

                if res.Msg.Name != req.Task.Name {
                        t.Errorf("Expected task name %s, got %s", req.Task.Name, res.Msg.Name)
                }
        })
}
