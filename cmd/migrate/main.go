package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/I-Van-Radkov/url-shortener/internal/config"
	"github.com/I-Van-Radkov/url-shortener/pkg/logger"
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

type MigrationConfig struct {
	Env      string `env:"ENV"`
	LogLevel string `env:"LOG_LEVEL"`

	Postgres config.PostgresConfig
}

func parseMigrationConfig() (*MigrationConfig, error) {
	cfg := &MigrationConfig{}

	_ = cleanenv.ReadConfig("config/.env", cfg)
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("parse migration config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("validate migration config: %w", err)
	}

	return cfg, nil
}

func (c *MigrationConfig) validate() error {
	if c.Env != "dev" && c.Env != "prod" {
		return fmt.Errorf("ENV must be dev or prod, got %q", c.Env)
	}

	if c.LogLevel != "debug" && c.LogLevel != "info" &&
		c.LogLevel != "warn" && c.LogLevel != "error" {
		return fmt.Errorf("LOG_LEVEL must be debug, info, warn or error, got %q", c.LogLevel)
	}

	if c.Postgres.Username == "" {
		return fmt.Errorf("POSTGRES_USER is required")
	}
	if c.Postgres.Password == "" {
		return fmt.Errorf("POSTGRES_PASSWORD is required")
	}
	if c.Postgres.Host == "" {
		return fmt.Errorf("POSTGRES_HOST is required")
	}
	if c.Postgres.Port == "" {
		return fmt.Errorf("POSTGRES_PORT is required")
	}
	if c.Postgres.DbName == "" {
		return fmt.Errorf("POSTGRES_DB is required")
	}

	return nil
}

func main() {
	cfg, err := parseMigrationConfig()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	lg, err := logger.NewLogger(cfg.Env, cfg.LogLevel)
	if err != nil {
		log.Fatalf("logger init error: %v", err)
	}
	defer lg.Sync()

	lg.Info("starting database migrations",
		zap.String("host", cfg.Postgres.Host),
		zap.String("db_name", cfg.Postgres.DbName),
	)

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Postgres.Username,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DbName,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		lg.Error("failed to open database connection", zap.Error(err))
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		lg.Error("failed to ping database", zap.Error(err))
		os.Exit(1)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		lg.Error("failed to set goose dialect", zap.Error(err))
		os.Exit(1)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		lg.Error("failed to apply migrations", zap.Error(err))
		os.Exit(1)
	}

	lg.Info("database migrations completed successfully")
}
