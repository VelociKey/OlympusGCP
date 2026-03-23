package inference

import (
	"context"
	"log/slog"

	firebasev1 "olympus.fleet/00SDLC/OlympusGCP/gen/firebase"
	"connectrpc.com/connect"
)

type FirebaseServer struct{}

func (s *FirebaseServer) CreateUser(ctx context.Context, req *connect.Request[firebasev1.CreateUserRequest]) (*connect.Response[firebasev1.CreateUserResponse], error) {
	slog.Info("CreateUser", "email", req.Msg.Email)
	return connect.NewResponse(&firebasev1.CreateUserResponse{Uid: "user-abc"}), nil
}

func (s *FirebaseServer) SetDocument(ctx context.Context, req *connect.Request[firebasev1.SetDocumentRequest]) (*connect.Response[firebasev1.SetDocumentResponse], error) {
	slog.Info("SetDocument", "collection", req.Msg.Collection, "id", req.Msg.DocId)
	return connect.NewResponse(&firebasev1.SetDocumentResponse{Uid: req.Msg.DocId}), nil
}
