package main

import (
	"fmt"
	"github.com/take-bit/auth-service/internal/config"
	"github.com/take-bit/auth-service/pkg/logging"

	"log/slog"
)

func main() {
	cfg := config.MustLoadConfig()

	logger := logging.NewLogger(cfg.Logging.Level)

	logger.Info("Starting application",
		slog.String("Logging level", cfg.Logging.Level),
		slog.String("Address server", fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)))

	logger.Info("test")
}
