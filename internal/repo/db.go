package repo

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type DB struct {
	Pool *pgxpool.Pool
}

func Connect(ctx context.Context) (*DB, error) {
	cfg, err := pgxpool.ParseConfig(viper.GetString("postgres.dsn"))
	if err != nil {
		return nil, err
	}
	cfg.MaxConns = int32(viper.GetInt("postgres.max_conns"))
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	log.Info().Msg("connected to postgres")
	return &DB{Pool: pool}, nil
}

func (d *DB) Close() {
	if d != nil && d.Pool != nil {
		d.Pool.Close()
	}
}
