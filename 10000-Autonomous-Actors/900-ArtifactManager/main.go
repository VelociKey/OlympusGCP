package main

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os/exec"
	"time"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	whisper "olympus.fleet/00SDLC/Olympus2/90000-Enablement-Labs/90200-Logic-Libraries/220-Whisper"
	artifactsv1 "olympus.fleet/00SDLC/OlympusGCP/gen/olympus/artifacts/v1"
	"olympus.fleet/00SDLC/OlympusGCP/gen/olympus/artifacts/v1/artifactsv1connect"
)

type ArtifactServer struct {
	logger *whisper.WhisperLog
}

func (s *ArtifactServer) PushImage(ctx context.Context, req *connect.Request[artifactsv1.PushImageRequest]) (*connect.Response[artifactsv1.PushImageResponse], error) {
	start := time.Now()
	slog.Info("ArtifactManager: Pushing OCI Image to Local Registry", "image", req.Msg.ImageName, "tag", req.Msg.Tag)

	// Substrate: Import TAR into Podman/Docker local storage
	cmd := exec.Command("podman", "load")
	cmd.Stdin = bytes.NewReader(req.Msg.Tarball)
	if err := cmd.Run(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to load image into substrate: %v", err))
	}

	s.logger.Log("PUSH_IMAGE", "SUCCESS", req.Msg.ImageName, "SUBSTRATE_LOADED", time.Since(start))
	return connect.NewResponse(&artifactsv1.PushImageResponse{
		Digest: fmt.Sprintf("sha256:emulated-%d", time.Now().Unix()),
		Status: "UPLOADED_TO_WORKSTATION_STORAGE",
	}), nil
}

func (s *ArtifactServer) ListImages(ctx context.Context, req *connect.Request[artifactsv1.ListImagesRequest]) (*connect.Response[artifactsv1.ListImagesResponse], error) {
	cmd := exec.Command("podman", "images", "--format", "{{.Repository}}")
	output, err := cmd.Output()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	lines := bytes.Split(bytes.TrimSpace(output), []byte("\n"))
	var names []string
	for _, l := range lines {
		names = append(names, string(l))
	}

	return connect.NewResponse(&artifactsv1.ListImagesResponse{ImageNames: names}), nil
}

func main() {
	w := whisper.New("ArtifactManager", "gcp_artifacts.lpsv")
	defer w.Close()

	server := &ArtifactServer{logger: w}
	mux := http.NewServeMux()
	path, handler := artifactsv1connect.NewArtifactServiceHandler(server)
	mux.Handle(path, handler)

	port := "8099"
	slog.Info("ArtifactManager: Booting OCI-Compatible High-Fidelity Substrate...", "port", port)

	http.ListenAndServe(
		"localhost:"+port,
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
