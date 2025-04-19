package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prathoss/integration_testing/service_picture/config"
	"github.com/prathoss/integration_testing/service_picture/repository"
)

type Server struct {
	config                config.Config
	pictureRepository     repository.Picture
	pictureViewRepository repository.View
	dbConnPool            *pgxpool.Pool
	// closers are dependencies, that should be closed when shutting down server, they are registered on server shutdown
	closers []closer
}

type closer interface {
	Close()
}

func New(cfg config.Config) (*Server, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pgConn, err := pgxpool.New(ctx, cfg.DatabaseURI)
	if err != nil {
		return nil, err
	}

	return &Server{
		config:                cfg,
		closers:               []closer{pgConn},
		pictureRepository:     repository.NewPicture(pgConn),
		pictureViewRepository: repository.NewView(pgConn),
		dbConnPool:            pgConn,
	}, nil
}

func (s *Server) Run(ctx context.Context) error {
	mux := http.NewServeMux()
	s.setupRoutes(mux)

	hs := http.Server{
		Addr:         fmt.Sprintf("%s:%d", s.config.ServerAddress, s.config.ServerPort),
		ReadTimeout:  s.config.ServerReadTimeout,
		WriteTimeout: s.config.ServerWriteTimeout,
		Handler:      mux,
		ErrorLog:     slog.NewLogLogger(slog.Default().Handler(), slog.LevelError),
	}

	for _, closer := range s.closers {
		hs.RegisterOnShutdown(
			func() {
				closer.Close()
			},
		)
	}

	errChan := make(chan error)
	go func() {
		slog.InfoContext(ctx, "starting server", slog.String("address", hs.Addr))
		if err := hs.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
			return
		}
		slog.InfoContext(ctx, "shutting down server")
		errChan <- nil
	}()

	var err error
	select {
	case <-ctx.Done():
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		err = hs.Shutdown(shutdownCtx)
	case err = <-errChan:
	}
	return err
}
