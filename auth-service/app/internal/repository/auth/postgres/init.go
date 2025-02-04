package postgresAuth

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"log/slog"
)

var (
	ErrUserNotSaved   = errors.New("failed to save user")
	ErrUsernameExists = errors.New("username already exists")
	ErrEmailExists    = errors.New("email already exists")
)

type Client interface {
	StartOperation() error
	EndOperation()
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}
type Repository struct {
	client Client
	logger *slog.Logger
}

func NewRepository(client Client, logger *slog.Logger) *Repository {
	return &Repository{
		client: client,
		logger: logger,
	}
}
