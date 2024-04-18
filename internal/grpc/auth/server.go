package auth

import (
	ssov1 "EMP_Back/gen/go/sso"
	"context"
	_ "github.com/go-ozzo/ozzo-validation"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	//TODO: push and import proto file
)

type serverApi struct {
	ssov1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	ssov1.RegisterAuthServer(gRPC, &serverApi{})
}

func (s *serverApi) Login(ctx context.Context,
	req *ssov1.LoginRequest,
) (
	*ssov1.LoginResponse, error) {

	err := validation.Validate(req.GetEmail(), validation.Required, is.Email, validation.Length(3, 100))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid email")
	}
	err = validation.Validate(req.GetPassword(), validation.Required, validation.Length(5, 20))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid password")
	}
	err = validation.Validate(req.GetUsername(), validation.Required, is.Alphanumeric, validation.Length(3, 20))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid username")
	}
	err = validation.Validate(req.GetAppId(), validation.Required, is.Alphanumeric, validation.Length(1, 10))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid app id")
	}

	return &ssov1.LoginResponse{
		Token: "",
	}, nil
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
