package main

import (
        "context"
        "fmt"
        "log/slog"
        "net/http"
        "os"
        "os/signal"
        "syscall"
        "time"

        "olympus.fleet/00SDLC/OlympusGCP/gen/google/cloud/tasks/v2/cloudtaskspbconnect"
        // taskspbconnect is used as alias above if needed, but let's check eventsv1 too
        "olympus.fleet/00SDLC/OlympusGCP/gen/olympus/tasks/v1/tasksv1connect"
        "olympus.fleet/00SDLC/OlympusGCP/10000-Autonomous-Actors/10700-Processing-Engines/10710-Reasoning-Inference/inference"

        "golang.org/x/net/http2"
        "golang.org/x/net/http2/h2c"
)

func main() {
        storageDir := "00SDLC/OlympusGCP/C0990-Ephemeral-Scratch"
        eventsServer := inference.NewEventsServer(storageDir)
        cloudTasksServer := inference.NewCloudTasksServer(eventsServer)
        
        mux := http.NewServeMux()
        
        // Legacy Events Service
        path, handler := eventsv1connect.NewEventsServiceHandler(eventsServer)
        mux.Handle(path, handler)

        // New High-Fidelity Cloud Tasks v2 Service
        tasksPath, tasksHandler := taskspbconnect.NewCloudTasksHandler(cloudTasksServer)
        mux.Handle(tasksPath, tasksHandler)

        // Health Check / Pulse
        mux.HandleFunc("/pulse", func(w http.ResponseWriter, r *http.Request) {
                w.Header().Set("Content-Type", "application/json")
                fmt.Fprintf(w, `{"status":"HEALTHY", "workspace":"OlympusGCP-Events", "time":"%s"}`, time.Now().Format(time.RFC3339))
        })

        port := "8094"
        slog.Info("EventsManager starting", "port", port)

        srv := &http.Server{
                Addr:              ":" + port,
                Handler:           h2c.NewHandler(mux, &http2.Server{}),
                ReadHeaderTimeout: 3 * time.Second,
        }

        done := make(chan os.Signal, 1)
        signal.Notify(done, os.Interrupt, syscall.SIGTERM)

        go func() {
                if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed { 
                        slog.Error("Server failed", "error", err)
                }
        }()

        <-done
        slog.Info("EventsManager shutting down...")
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        srv.Shutdown(ctx)
}
