package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/I-Van-Radkov/url-shortener/internal/config"
	"github.com/I-Van-Radkov/url-shortener/pkg/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.ParseConfigFromEnv()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	lg, err := logger.NewLogger(cfg.Env, cfg.LogLevel)
	if err != nil {
		log.Fatalf("logger init error: %v", err)
	}
	defer lg.Sync()

	lg.Info("starting database migrations",
		zap.String("host", cfg.PostgresConfig.Host),
		zap.String("db_name", cfg.PostgresConfig.DbName),
	)

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.PostgresConfig.Username,
		cfg.PostgresConfig.Password,
		cfg.PostgresConfig.Host,
		cfg.PostgresConfig.Port,
		cfg.PostgresConfig.DbName,
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
