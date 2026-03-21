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

	"olympus.fleet/00SDLC/OlympusGCP-Intelligence/40000-Communication-Contracts/40400-Protocol-Synthetics/connect-rpc/gen/v1/intelligence/intelligencev1connect"
	"olympus.fleet/00SDLC/OlympusGCP-Intelligence/10000-Autonomous-Actors/10700-Processing-Engines/10710-Reasoning-Inference/inference"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	server := &inference.IntelligenceServer{}
	mux := http.NewServeMux()
	path, handler := intelligencev1connect.NewIntelligenceServiceHandler(server)
	mux.Handle(path, handler)

	// Health Check / Pulse
	mux.HandleFunc("/pulse", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"HEALTHY", "workspace":"OlympusGCP-Intelligence", "time":"%s"}`, time.Now().Format(time.RFC3339))
	})

	port := "8096" // From genesis.json
	slog.Info("IntelligenceManager starting", "port", port)

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
	slog.Info("IntelligenceManager shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
