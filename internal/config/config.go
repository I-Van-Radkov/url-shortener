package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type PostgresConfig struct {
	Username string `env:"POSTGRES_USER" env-default:"postgres"`
	Password string `env:"POSTGRES_PASSWORD" env-default:"postgres"`
	Host     string `env:"POSTGRES_HOST" env-default:"db"`
	Port     string `env:"POSTGRES_PORT" env-default:"5432"`
	DbName   string `env:"POSTGRES_DB" env-default:"postgres"`
}

type Config struct {
	GracefulShutdownTimeout time.Duration `env:"GRACEFUL_SHUTDOWN_TIMEOUT" env-default:"5s"`
	StorageType             string        `env:"STORAGE_TYPE" env-default:"memory"`
	MaxAttemptsToGen        int           `env:"MAX_ATTEMPTS_TO_GEN" env-default:"5"`

	Port         int           `env:"PORT" env-default:"8080"`
	ReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT" env-default:"30s"`
	WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT" env-default:"30s"`

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

	return nil
}

func ParseConfigFromEnv() (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config from env: %w", err)
	}

	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate config: %w", err)
	}

	return cfg, nil
}
