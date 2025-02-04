package grpcAuthController

import (
	"context"
	"github.com/take-bit/auth-service/internal/domain/auth/dto"
	grpcGen "github.com/take-bit/proto/gen/auth/go/v1"
	"log/slog"
)

type RegisterUser interface {
	RegisterUser(ctx context.Context, req *dto.RegisterUserRequest) (*dto.RegisterUserResponse, error)
}

type LoginUser interface {
	Execute(ctx context.Context, req *dto.LoginUserRequest) (*dto.LoginUserResponse, error)
}

type UpdateAccessToken interface {
	Execute(ctx context.Context, req *dto.UpdateAccessTokenRequest) (*dto.UpdateAccessTokenResponse, error)
}

type Controller struct {
	logger *slog.Logger
	grpcGen.UnimplementedAuthServiceServer
	registerUser RegisterUser
	loginUser    LoginUser
	updateToken  UpdateAccessToken
}

func NewController(logger *slog.Logger, registerUser RegisterUser, loginUser LoginUser, token UpdateAccessToken) *Controller {
	return &Controller{
		registerUser: registerUser,
		loginUser:    loginUser,
		updateToken:  token,
		logger:       logger,
	}
}
