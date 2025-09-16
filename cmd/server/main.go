package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"go-task-tracker/internal/api"
	"go-task-tracker/internal/repo"
)

func main() {
	configureLogging()
	if err := loadConfig(); err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	addr := fmt.Sprintf("%s:%d", viper.GetString("server.host"), viper.GetInt("server.port"))

	ctx := context.Background()

	db, err := repo.Connect(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("db connect failed")
	}
	defer db.Close()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(requestLogger())

	api.NewAPI(db).RegisterRoutes(r)

	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Info().Str("addr", addr).Msg("http server starting")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("http server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Error().Err(err).Msg("server shutdown error")
	}
	log.Info().Msg("server stopped")
}

func configureLogging() {
	if os.Getenv("ENV") == "dev" || os.Getenv("ENV") == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		return
	}

	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func loadConfig() error {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("./configs")
	viper.AutomaticEnv()

	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("jwt.secret", "change-me")
	viper.SetDefault("jwt.access_ttl", "15m")
	viper.SetDefault("postgres.dsn", "postgres://postgres:postgres@localhost:5432/taskdb?sslmode=disable")
	viper.SetDefault("postgres.max_conns", 10)

	if err := viper.ReadInConfig(); err != nil {
		log.Warn().Err(err).Msg("using defaults/env; config file not loaded")
	}
	return nil
}

func requestLogger() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			log.Info().Str("method", r.Method).Str("path", r.URL.Path).Dur("duration", time.Since(start)).Msg("request")
		})
	}
}
