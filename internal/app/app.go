package app

import (
	"fmt"
	"go-task-tracker/config"
	v1 "go-task-tracker/internal/api/v1"
	"go-task-tracker/internal/repo"
	service "go-task-tracker/internal/services"
	"go-task-tracker/pkg/hasher"
	httpserver "go-task-tracker/pkg/httpsserver"
	"go-task-tracker/pkg/postgres"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func Run(configPath string) {
	// Configuration
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatal().Err(err).Msg("config error")
	}

	// Logger
	configureLogging()

	// DB
	log.Info().Msg("Initializing postgres...")
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.MaxPoolSize))
	if err != nil {
		log.Err(err).Msg("app - Run - pgdb.NewServices")
	}
	defer pg.Close()

	// Repositories
	log.Info().Msg("Initializing repositories...")
	repositories := repo.NewRepositories(pg)

	// Services dependencies
	log.Info().Msg("Initializing services...")
	deps := service.ServicesDependencies{
		Repos: repositories,
		// GDrive:   gdrive.New(cfg.WebAPI.GDriveJSONFilePath),
		Hasher:   hasher.NewSHA1Hasher(cfg.Hasher.Salt),
		SignKey:  cfg.JWT.SignKey,
		TokenTTL: cfg.JWT.TokenTTL,
	}
	services := service.NewServices(deps)

	// Handlers
	log.Info().Msg("Initializing handlers and routes...")

	handler := chi.NewRouter()
	v1.RegisterRoutes(handler, services)

	// HTTP server
	log.Info().Msg("Starting http server...")
	log.Debug().Msg(fmt.Sprintf("Server port: %s", cfg.HTTP.Port))
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	log.Info().Msg("Configuring graceful shutdown...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info().Msg("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Error().Err(err).Msg("app - Run - httpServer.Notify")
	}

	// Graceful shutdown
	log.Info().Msg("Shutting down...")
	err = httpServer.Shutdown()
	if err != nil {
		log.Error().Err(err).Msg("app - Run - httpServer.Shutdown")
	}
}
