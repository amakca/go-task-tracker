package app

import (
	"context"
	"go-task-tracker/config"
	v1 "go-task-tracker/internal/api/v1"
	"go-task-tracker/internal/repo"
	service "go-task-tracker/internal/services"
	"go-task-tracker/pkg/hasher"
	httpserver "go-task-tracker/pkg/httpsserver"
	"go-task-tracker/pkg/logctx"
	"go-task-tracker/pkg/postgres"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
)

func Run(configPath string) {
	// Configuration
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		slog.Error("config error", "err", err)
		os.Exit(1)
	}

	// Logger
	configureLogging()
	// root context logger
	ctx := context.Background()
	ctx = logctx.WithLogger(ctx, slog.With("app", "go-task-tracker"))
	log := logctx.FromContext(ctx)

	// DB
	log.Info("Initializing postgres...")
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.MaxPoolSize))
	if err != nil {
		log.Error("app - Run - pgdb.NewServices", "err", err)
	}
	defer pg.Close()

	// Repositories
	log.Info("Initializing repositories...")
	repositories := repo.NewRepositories(pg)

	// Services dependencies
	log.Info("Initializing services...")
	deps := service.ServicesDependencies{
		Repos: repositories,
		// GDrive:   gdrive.New(cfg.WebAPI.GDriveJSONFilePath),
		Hasher:   hasher.NewSHA1Hasher(cfg.Hasher.Salt),
		SignKey:  cfg.JWT.SignKey,
		TokenTTL: cfg.JWT.TokenTTL,
	}
	services := service.NewServices(deps)

	// Handlers
	log.Info("Initializing handlers and routes...")
	r := chi.NewRouter()
	v1.NewRouter(r, services)

	// HTTP server
	log.Info("Starting http server...")
	log.Debug("Server starting", "port", cfg.HTTP.Port)
	httpServer := httpserver.New(r, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	log.Info("Configuring graceful shutdown...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal", "signal", s.String())
	case err = <-httpServer.Notify():
		log.Error("app - Run - httpServer.Notify", "err", err)
	}

	// Graceful shutdown
	log.Info("Shutting down...")
	err = httpServer.Shutdown()
	if err != nil {
		log.Error("app - Run - httpServer.Shutdown", "err", err)
	}
}
