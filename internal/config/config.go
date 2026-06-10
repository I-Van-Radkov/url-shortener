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
	Port                    int           `env:"PORT" env-default:"8080"`
	ReadTimeout             time.Duration `env:"HTTP_READ_TIMEOUT" env-default:"30s"`
	WriteTimeout            time.Duration `env:"HTTP_WRITE_TIMEOUT" env-default:"30s"`

	PostgresConfig PostgresConfig
}

func ParseConfigFromEnv() (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config from env: %w", err)
	}

	return cfg, nil
}
