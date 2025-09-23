package app

import (
	"go-task-tracker/config"
	v1 "go-task-tracker/internal/api/v1"
	"go-task-tracker/internal/repo"
	service "go-task-tracker/internal/services"
	"go-task-tracker/pkg/hasher"
	httpserver "go-task-tracker/pkg/httpsserver"
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

	// DB
	slog.Info("Initializing postgres...")
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.MaxPoolSize))
	if err != nil {
		slog.Error("app - Run - pgdb.NewServices", "err", err)
	}
	defer pg.Close()

	// Repositories
	slog.Info("Initializing repositories...")
	repositories := repo.NewRepositories(pg)

	// Services dependencies
	slog.Info("Initializing services...")
	deps := service.ServicesDependencies{
		Repos: repositories,
		// GDrive:   gdrive.New(cfg.WebAPI.GDriveJSONFilePath),
		Hasher:   hasher.NewSHA1Hasher(cfg.Hasher.Salt),
		SignKey:  cfg.JWT.SignKey,
		TokenTTL: cfg.JWT.TokenTTL,
	}
	services := service.NewServices(deps)

	// Handlers
	slog.Info("Initializing handlers and routes...")
	r := chi.NewRouter()
	v1.NewRouter(r, services)

	// HTTP server
	slog.Info("Starting http server...")
	slog.Debug("Server starting", "port", cfg.HTTP.Port)
	httpServer := httpserver.New(r, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	slog.Info("Configuring graceful shutdown...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		slog.Info("app - Run - signal", "signal", s.String())
	case err = <-httpServer.Notify():
		slog.Error("app - Run - httpServer.Notify", "err", err)
	}

	// Graceful shutdown
	slog.Info("Shutting down...")
	err = httpServer.Shutdown()
	if err != nil {
		slog.Error("app - Run - httpServer.Shutdown", "err", err)
	}
}
