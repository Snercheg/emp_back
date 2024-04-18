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

type Auth interface {
	Login(ctx context.Context,
		username string,
		email string,
		password string,
		appID int,
	) (token string, err error)
	RegisterNewUser(ctx context.Context,
		username string,
		email string,
		password string,
	) (userID int64, err error)
	IsAdmin(ctx context.Context,
		userID int64,
	) (isAdmin bool, err error)
}

type serverApi struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverApi{auth: auth})
}

func (s *serverApi) Login(ctx context.Context,
	req *ssov1.LoginRequest,
) (
	*ssov1.LoginResponse, error) {

	if err := validateLogin(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetUsername(), req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		// TODO: handle error type
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverApi) Register(
	ctx context.Context,
	req *ssov1.RegisterRequest,
) (
	*ssov1.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}
	userID, err := s.auth.RegisterNewUser(ctx, req.GetUsername(), req.GetEmail(), req.GetPassword())
	if err != nil {
		// TODO: handle error type
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &ssov1.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *serverApi) isAdmin(
	ctx context.Context,
	req *ssov1.IsAdminRequest) (
	*ssov1.IsAdminResponse, error) {

	if err := validateIsAdmin(req); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		// TODO: handle error type
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

func validateLogin(req *ssov1.LoginRequest) error {
	err := validation.Validate(req.GetEmail(), validation.Required, is.Email, validation.Length(3, 100))
	if err != nil {
		return status.Error(codes.InvalidArgument, "invalid email")
	}
	err = validation.Validate(req.GetPassword(), validation.Required, validation.Length(5, 20))
	if err != nil {
		return status.Error(codes.InvalidArgument, "invalid password")
	}
	err = validation.Validate(req.GetUsername(), validation.Required, is.Alphanumeric, validation.Length(3, 20))
	if err != nil {
		return status.Error(codes.InvalidArgument, "invalid username")
	}
	err = validation.Validate(req.GetAppId(), validation.Required, is.Alphanumeric, validation.Length(1, 10))
	if err != nil {
		return status.Error(codes.InvalidArgument, "invalid app id")
	}
	return nil
}

func validateRegister(req *ssov1.RegisterRequest) error {
	err := validation.Validate(req.GetEmail(), validation.Required, is.Email, validation.Length(3, 100))
	if err != nil {
		return status.Error(codes.InvalidArgument, "invalid email")
	}
	err = validation.Validate(req.GetPassword(), validation.Required, validation.Length(5, 20))
	if err != nil {
		return status.Error(codes.InvalidArgument, "invalid password")
	}
	err = validation.Validate(req.GetUsername(), validation.Required, is.Alphanumeric, validation.Length(3, 20))
	if err != nil {
		return status.Error(codes.InvalidArgument, "invalid username")
	}
	return nil
}

func validateIsAdmin(req *ssov1.IsAdminRequest) error {
	if err := validation.Validate(req.GetUserId(), validation.Required, is.Alphanumeric, validation.Length(1, 20)); err != nil {
		return status.Error(codes.InvalidArgument, "invalid user id")
	}
	return nil
}
