package main

import (
	"log"
	"os"

	"github.com/I-Van-Radkov/url-shortener/internal/app"
	"github.com/I-Van-Radkov/url-shortener/internal/config"
	"github.com/I-Van-Radkov/url-shortener/pkg/logger"
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

	lg.Info("logger initialized",
		zap.String("env", cfg.Env),
		zap.String("log_level", cfg.LogLevel),
	)

	a, err := app.NewApp(cfg, lg)
	if err != nil {
		lg.Error("failed to initialize application", zap.Error(err))
		os.Exit(1)
	}

	if err = a.Run(); err != nil {
		lg.Error("application stopped with error", zap.Error(err))
		os.Exit(1)
	}
}
