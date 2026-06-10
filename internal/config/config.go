package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type PostgresConfig struct {
	Username string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
	DbName   string `env:"POSTGRES_DB"`
}

type Config struct {
	Env      string `env:"ENV"`
	LogLevel string `env:"LOG_LEVEL"`

	GracefulShutdownTimeout time.Duration `env:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	StorageType             string        `env:"STORAGE_TYPE"`
	MaxAttemptsToGen        int           `env:"MAX_ATTEMPTS_TO_GEN"`

	Port         int           `env:"PORT"`
	ReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT"`
	WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT"`

	PostgresConfig PostgresConfig
}

func (c *Config) Validate() error {
	if c.StorageType != "memory" && c.StorageType != "postgres" {
		return fmt.Errorf("invalid STORAGE_TYPE: %s", c.StorageType)
	}

	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("invalid PORT: %d", c.Port)
	}

	if c.MaxAttemptsToGen < 1 {
		return fmt.Errorf("MAX_ATTEMPTS_TO_GEN must be between 1 and 10, got %d", c.MaxAttemptsToGen)
	}

	if c.GracefulShutdownTimeout < 1*time.Second {
		return fmt.Errorf("GRACEFUL_SHUTDOWN_TIMEOUT must be at least 1s, got %v", c.GracefulShutdownTimeout)
	}

	if c.Env != "dev" && c.Env != "prod" {
		return fmt.Errorf("invalid ENV: %s", c.Env)
	}

	if c.LogLevel != "debug" && c.LogLevel != "info" &&
		c.LogLevel != "warn" && c.LogLevel != "error" {
		return fmt.Errorf("invalid LOG_LEVEL: %s", c.LogLevel)
	}

	return nil
}

func ParseConfigFromEnv() (*Config, error) {
	cfg := &Config{}

	if _, err := os.Stat("config/.env"); err == nil {
		if err := cleanenv.ReadConfig("config/.env", cfg); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("failed to stat config file: %w", err)
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config from env: %w", err)
	}

	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate config: %w", err)
	}

	return cfg, nil
}
