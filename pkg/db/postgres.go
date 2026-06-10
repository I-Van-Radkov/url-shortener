package db

import (
	"context"
	"fmt"

	"github.com/I-Van-Radkov/url-shortener/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool
}

func NewPostgres(config config.PostgresConfig) (*Database, error) {
	dataSource := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		config.Username, config.Password, config.Host, config.Port, config.DbName)

	pool, err := pgxpool.New(context.Background(), dataSource)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return &Database{
		Pool: pool,
	}, nil
}

func (d *Database) Close() {
	if d.Pool != nil {
		d.Pool.Close()
	}
}
