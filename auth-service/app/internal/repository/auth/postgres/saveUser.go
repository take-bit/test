package postgresAuth

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"log/slog"
)

const (
	querySaveUser = "INSERT INTO users (uuid, username, email, hash_password) VALUES ($1, $2, $3, $4)"
)

func (r *Repository) SaveUser(ctx context.Context, uuid string, username string, email string, hashPassword string) error {
	const op = "repository/postgres/auth/SaveUser"

	requestID := ctx.Value("request").(string)

	logger := r.logger.With(
		slog.String("operation", op),
		slog.String("request_id", requestID))

	err := r.client.StartOperation()
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	defer r.client.EndOperation()

	rows, err := r.client.Exec(ctx, querySaveUser, uuid, username, email, hashPassword)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				if pgErr.ConstraintName == "users_username_key" {
					logger.Warn("Username already exists", slog.String("error", err.Error()))
					return fmt.Errorf("%s: %w", op, ErrUsernameExists)
				}
				if pgErr.ConstraintName == "users_email_key" {
					logger.Warn("Email already exists", slog.String("error", err.Error()))
					return fmt.Errorf("%s: %w", op, ErrEmailExists)
				}
			}
		}

		logger.Error("Failed to execute query", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, ErrUserNotSaved)
	}

	if rows.RowsAffected() == 0 {
		logger.Warn("No rows affected, user not saved")
		return ErrUserNotSaved
	}

	logger.Info("User successfully saved")
	return nil
}
