package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"olympus.fleet/00SDLC/OlympusGCP-Observability/40000-Communication-Contracts/40400-Protocol-Synthetics/connect-rpc/gen/v1/observability/observabilityv1connect"
	"olympus.fleet/00SDLC/OlympusGCP-Observability/10000-Autonomous-Actors/10700-Processing-Engines/10710-Reasoning-Inference/inference"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	server := &inference.ObservabilityServer{}
	mux := http.NewServeMux()
	path, handler := observabilityv1connect.NewObservabilityServiceHandler(server)
	mux.Handle(path, handler)

	// Health Check / Pulse
	mux.HandleFunc("/pulse", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"HEALTHY", "workspace":"OlympusGCP-Observability", "time":"%s"}`, time.Now().Format(time.RFC3339))
	})

	port := "8097" // From genesis.json
	slog.Info("ObservabilityManager starting", "port", port)

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           h2c.NewHandler(mux, &http2.Server{}),
		ReadHeaderTimeout: 3 * time.Second,
	}
	err := srv.ListenAndServe()
	if err != nil {
		slog.Error("Server failed", "error", err)
	}
}
