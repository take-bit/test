package postgresClient

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"sync/atomic"
	"time"
)

var (
	ErrDatabaseClosed = errors.New("the database does not accept requests")
)

type Client struct {
	pool   *pgxpool.Pool
	closed atomic.Int32
	count  atomic.Int32
	logger *slog.Logger
}

func NewClient(ctx context.Context, connStr string, maxAttempt int, logger *slog.Logger) (*Client, error) {
	const op = "postgresClient/NewClient"

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		logger.Error("Failed to parse config",
			slog.String("operation", op),
			slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: failed to parse config: %w", op, err)
	}

	for i := 0; i < maxAttempt; i++ {
		pool, err := pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			logger.Error("Failed to create connection pool",
				slog.String("operation", op),
				slog.String("attempt", fmt.Sprintf("%d", i+1)),
				slog.String("error", err.Error()))
			time.Sleep(time.Duration(2<<i) * time.Second)
			continue
		}

		if err := pool.Ping(ctx); err != nil {
			logger.Error("Failed to ping database",
				slog.String("operation", op),
				slog.String("attempt", fmt.Sprintf("%d", i+1)),
				slog.String("error", err.Error()))
			time.Sleep(time.Duration(2<<i) * time.Second)
			continue
		}

		logger.Info("Successfully connected to the database",
			slog.String("operation", op))
		return &Client{
			pool:   pool,
			logger: logger,
		}, nil
	}

	logger.Error("Failed to establish connection to the database after maximum attempts",
		slog.String("operation", op),
		slog.Int("max_attempts", maxAttempt))
	return nil, fmt.Errorf("%s: the connection to the database could not be established", op)
}

func (c *Client) StartOperation() error {
	if c.closed.Load() > 0 {
		c.logger.Warn("Attempted operation on a closed database",
			slog.String("operation", "StartOperation"))
		return ErrDatabaseClosed
	}

	c.count.Add(1)
	c.logger.Debug("Operation started",
		slog.String("operation", "StartOperation"),
		slog.Int("active_operations", int(c.count.Load())))

	return nil
}

func (c *Client) EndOperation() {
	old := c.count.Load()
	c.count.Store(old - 1)
	c.logger.Debug("Operation ended",
		slog.String("operation", "EndOperation"),
		slog.Int("active_operations", int(c.count.Load())))
}

func (c *Client) Close() {
	if c.closed.CompareAndSwap(0, 1) {
		c.logger.Info("Starting to close the database connection pool")
		for c.count.Load() > 0 {
			c.logger.Debug("Waiting for active operations to finish before closing the database")
			time.Sleep(200 * time.Millisecond)
		}

		c.pool.Close()
		c.logger.Info("The database was successfully stopped")
	} else {
		c.logger.Warn("Attempted to close the database that was already closed",
			slog.String("operation", "Close"))
	}
}

func (c *Client) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return c.pool.Exec(ctx, sql, arguments...)
}
func (c *Client) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return c.pool.Query(ctx, sql, args...)
}
func (c *Client) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return c.pool.QueryRow(ctx, sql, args...)
}
