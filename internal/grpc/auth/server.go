package auth

import (
	"context"
	"google.golang.org/grpc"
	//TODO: push and import proto file
)

type serverApi struct {
	ssov1.UnumplimentedAuthServer
}

func Register(gRPC *grpc.Server) {
	ssov1.RegisterAuthServer(gRPC, &serverApi{})
}

func (s *serverApi) Login(ctx context.Context,
	req *ssov1.AuthRequest,
) (
	*ssov1.LoginResponse, error) {
	// TODO: Generate token
	panic("implement me")
}

func (s *serverApi) Register(
	ctx context.Context,
	req *ssov1.RegisterRequest,
) (
	*ssov1.RegisterResponse, error) {
	panic("implement me")
}

func (s *serverApi) isAdmin(
	ctx context.Context,
	req *ssov1.IsAdminRequest) (
	*ssov1.IsAdminResponse, error) {
	panic("implement me")
}
