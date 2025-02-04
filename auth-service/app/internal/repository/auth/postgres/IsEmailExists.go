package postgresAuth

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log/slog"
)

const (
	queryIsExistsEmail = "SELECT uuid FROM users WHERE email = $1"
)

func (r *Repository) IsEmailExists(ctx context.Context, email string) error {
	const op = "repository/postgres/auth/IsEmailExists"

	requestID := ctx.Value("request").(string)

	logger := r.logger.With(
		slog.String("operation", op),
		slog.String("request_id", requestID))

	err := r.client.StartOperation()
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	defer r.client.EndOperation()

	var uuid string
	err = r.client.QueryRow(ctx, queryIsExistsEmail, email).Scan(&uuid)

	if errors.Is(err, pgx.ErrNoRows) {
		logger.Info("Username is unique")
		return nil
	}

	if err != nil {
		logger.Error("Error executing query", slog.String("error", err.Error()))
		return fmt.Errorf("%s: query execution failed: %w", op, err)
	}

	logger.Warn("Username already exists")
	return fmt.Errorf("%s: %w", op, ErrUsernameExists)
}
