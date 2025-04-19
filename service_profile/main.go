package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/prathoss/integration_testing/logging"
	"github.com/prathoss/integration_testing/service_profile/config"
	"github.com/prathoss/integration_testing/service_profile/server"
)

func main() {
	err := func() error {
		ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
		defer cancel()

		cfg, err := config.NewFromEnv()
		if err != nil {
			slog.Error("failed to load configuration", logging.Err(err))
			return err
		}

		logging.Setup(cfg.LogLevel)

		s, err := server.New(cfg)
		if err != nil {
			slog.Error("failed to create server", logging.Err(err))
			return err
		}

		if err := s.Run(ctx); err != nil {
			slog.Error("problem running the server", logging.Err(err))
			return err
		}
		return nil
	}()
	if err != nil {
		os.Exit(1)
	}
}
