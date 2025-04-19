package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/prathoss/integration_testing/logging"
	"github.com/prathoss/integration_testing/service_picture/config"
	"github.com/prathoss/integration_testing/service_picture/server"
)

func main() {
	err := func() error {
		ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
		defer cancel()

		cfg, err := config.NewFromEnv()
		if err != nil {
			slog.Error("failed to load configuration", slog.Any("error", err))
			return err
		}

		// Set up logging
		logging.Setup(cfg.LogLevel)

		// Create and run server
		s, err := server.New(cfg)
		if err != nil {
			slog.Error("failed to create server", slog.Any("error", err))
			return err
		}

		if err := s.Run(ctx); err != nil {
			slog.Error("problem running the server", slog.Any("error", err))
			return err
		}
		return nil
	}()
	if err != nil {
		os.Exit(1)
	}
}
