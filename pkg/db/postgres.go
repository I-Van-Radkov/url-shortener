package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresConfig struct {
	Username string `env:"POSTGRES_USER" env-default:"postgres"`
	Password string `env:"POSTGRES_PASSWORD" env-default:"postgres"`
	Host     string `env:"POSTGRES_HOST" env-default:"db"`
	Port     string `env:"POSTGRES_PORT" env-default:"5432"`
	DbName   string `env:"POSTGRES_DB" env-default:"postgres"`
}

type Database struct {
	Pool *pgxpool.Pool
}

func New(config PostgresConfig) (*Database, error) {
	dataSource := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		config.Username, config.Password, config.Host, config.Port, config.DbName)

	pool, err := pgxpool.New(context.Background(), dataSource)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db after retries: %w", err)
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
