package authUseCase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/take-bit/auth-service/internal/domain/auth/dto"
	"github.com/take-bit/auth-service/pkg/security"
	"log/slog"
	"time"
)

type Repository interface {
	SaveUser(ctx context.Context, uuid string, username string, email string, hashPassword string) error
	SetRefreshToken(ctx context.Context, uuid string, refreshToken string) error
	IsUsernameExists(ctx context.Context, username string) error
	IsEmailExists(ctx context.Context, email string) error
}

type Service interface {
	HashPassword(password string) (string, error)
	VerificationEmail(email string) error
}

type RegisterUserUseCase struct {
	service         Service
	repo            Repository
	logger          *slog.Logger
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewRegisterUserUseCase(service Service, repo Repository, logger *slog.Logger) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		service: service,
		repo:    repo,
		logger:  logger,
	}
}

func (r *RegisterUserUseCase) RegisterUser(ctx context.Context, req *dto.RegisterUserRequest) (*dto.RegisterUserResponse, error) {
	const op = "useCase/auth/RegisterUser"

	requestID, ok := ctx.Value("request_id").(string)
	if !ok {
		requestID = "unknown"
	}

	logger := r.logger.With(
		slog.String("operation", op),
		slog.String("request_id", requestID))

	logger.Debug("Starting user registration",
		slog.String("username", req.Username),
		slog.String("email", req.Email))

	if err := r.repo.IsUsernameExists(ctx, req.Username); err != nil {
		logger.Warn("Username already exists", slog.String("username", req.Username))
		return nil, err
	}

	if err := r.repo.IsEmailExists(ctx, req.Email); err != nil {
		logger.Warn("Email already exists", slog.String("email", req.Email))
		return nil, err
	}

	hashPassword, err := r.service.HashPassword(req.Password)
	if err != nil {
		logger.Error("Failed to hash password", slog.String("error", err.Error()))
		return nil, err
	}

	userUUID := uuid.NewString()
	accessJIT := uuid.NewString()
	refreshJIT := uuid.NewString()

	if err := r.repo.SaveUser(ctx, userUUID, req.Username, req.Email, hashPassword); err != nil {
		logger.Error("Failed to save user", slog.String("error", err.Error()))
		return nil, err
	}
	logger.Debug("User saved to database", slog.String("uuid", userUUID))

	accessToken, err := security.CreateToken(userUUID, accessJIT, req.Username, false, r.AccessTokenTTL)
	if err != nil {
		logger.Error("Failed to create access token", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: failed to create access token: %w", op, err)
	}
	logger.Debug("Access token generated successfully")

	refreshToken, err := security.CreateToken(userUUID, refreshJIT, req.Username, false, r.RefreshTokenTTL)
	if err != nil {
		logger.Error("Failed to create refresh token", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: failed to create refresh token: %w", op, err)
	}
	logger.Debug("Refresh token generated successfully")

	if err := r.repo.SetRefreshToken(ctx, userUUID, refreshToken); err != nil {
		logger.Error("Failed to save refresh token", slog.String("error", err.Error()))
		return nil, err
	}
	logger.Debug("Refresh token saved to database")

	go func(email string) {
		logger.Debug("Starting email verification process", slog.String("email", email))
		if err := r.service.VerificationEmail(email); err != nil {
			logger.Error("Failed to send verification email", slog.String("error", err.Error()))
		} else {
			logger.Debug("Verification email sent successfully")
		}
	}(req.Email)

	logger.Info("User registered successfully", slog.String("uuid", userUUID))

	return &dto.RegisterUserResponse{
		UUID:         userUUID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
