package logger

import (
	"fmt"

	"go.uber.org/zap"
)

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Sync() error
}

type L struct {
	z *zap.Logger
}

func NewLogger(env, level string) (Logger, error) {
	var cfg zap.Config
	switch env {
	case "dev":
		cfg = zap.NewDevelopmentConfig()
	case "prod":
		cfg = zap.NewProductionConfig()
	default:
		return nil, fmt.Errorf("invalid env: %s", env)
	}

	switch level {
	case "debug":
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		cfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		return nil, fmt.Errorf("invalid level: %s", level)
	}

	z, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return &L{z: z}, nil
}

func (l *L) Debug(msg string, fields ...zap.Field) {
	l.z.Debug(msg, fields...)
}

func (l *L) Info(msg string, fields ...zap.Field) {
	l.z.Info(msg, fields...)
}

func (l *L) Warn(msg string, fields ...zap.Field) {
	l.z.Warn(msg, fields...)
}

func (l *L) Error(msg string, fields ...zap.Field) {
	l.z.Error(msg, fields...)
}

func (l *L) Sync() error {
	return l.z.Sync()
}
