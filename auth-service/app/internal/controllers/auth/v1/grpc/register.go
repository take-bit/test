package grpcAuthController

import (
	"context"
	"github.com/take-bit/auth-service/internal/domain/auth/dto"
	grpcGen "github.com/take-bit/proto/gen/auth/go/v1"
	"log/slog"
)

func (c *Controller) Register(ctx context.Context, req *grpcGen.RegisterRequest) (*grpcGen.RegisterResponse, error) {
	const op = "controllers/auth/v1/grpc/Register"

	requestID, ok := ctx.Value("request_id").(string)
	if !ok {
		requestID = "unknown"
	}

	logger := c.logger.With(
		slog.String("operation", op),
		slog.String("request_id", requestID))

	logger.Debug("Getting started with the handler",
		slog.String("username", req.GetUsername()),
		slog.String("email", req.GetEmail()),
		slog.String("password", req.GetPassword()))

	request := dto.RegisterUserRequest{
		Username: req.GetUsername(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	response, err := c.registerUser.Execute(ctx, &request)
	if err != nil {
		// todo

		return nil, nil
	}

	res := &grpcGen.RegisterResponse{Result: &grpcGen.RegisterResponse_Tokens{
		Tokens: &grpcGen.Tokens{
			AccessToken:  response.AccessToken,
			RefreshToken: response.RefreshToken,
		},
	}}

	return res, nil
}
