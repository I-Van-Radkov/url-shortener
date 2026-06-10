package app

import (
	"context"
	"errors"
	"fmt"
	nethttp "net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/I-Van-Radkov/url-shortener/internal/adapter/memory"
	"github.com/I-Van-Radkov/url-shortener/internal/adapter/postgres"
	"github.com/I-Van-Radkov/url-shortener/internal/config"
	httpserver "github.com/I-Van-Radkov/url-shortener/internal/controller/http"
	v1 "github.com/I-Van-Radkov/url-shortener/internal/controller/http/v1"
	"github.com/I-Van-Radkov/url-shortener/internal/model"
	"github.com/I-Van-Radkov/url-shortener/internal/usecase"
	"github.com/I-Van-Radkov/url-shortener/pkg/db"
	"github.com/I-Van-Radkov/url-shortener/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type App struct {
	httpServer *httpserver.Server
	cfg        *config.Config
	db         *db.Database
	logger     logger.Logger
}

func NewApp(cfg *config.Config, logger logger.Logger) (*App, error) {
	a := &App{
		cfg:    cfg,
		logger: logger,
	}

	a.logger.Info("initializing application",
		zap.String("storage_type", cfg.StorageType),
		zap.Int("port", cfg.Port),
	)

	gen := usecase.NewCodeGenerator()

	var uc *usecase.Usecase

	switch a.cfg.StorageType {
	case "memory":
		a.logger.Info("using in-memory storage")

		memRepo := memory.NewRepo()
		uc = usecase.NewUsecase(memRepo, gen, a.cfg.MaxAttemptsToGen)
	case "postgres":
		a.logger.Info("connecting to postgres",
			zap.String("host", cfg.PostgresConfig.Host),
			zap.String("db_name", cfg.PostgresConfig.DbName),
		)

		database, err := db.NewPostgres(cfg.PostgresConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}

		a.db = database
		a.logger.Info("postgres connection established")

		postgresRepo := postgres.NewRepo(database.Pool)
		uc = usecase.NewUsecase(postgresRepo, gen, a.cfg.MaxAttemptsToGen)
	default:
		return nil, model.ErrInvalidStorageType
	}

	handlerV1 := v1.NewHandler(uc)

	router := gin.Default()
	handlerV1.RegisterRoutesGin(router)

	httpServer := httpserver.NewServer(a.cfg.Port, a.cfg.ReadTimeout, a.cfg.WriteTimeout, router)
	a.httpServer = httpServer

	a.logger.Info("application initialized successfully")

	return a, nil
}

func (a *App) Run() error {
	if err := a.run(); err != nil {
		return err
	}

	return nil
}

func (a *App) run() error {
	defer func() {
		if a.db != nil {
			a.logger.Info("closing database connection pool")
			a.db.Close()
		}

		_ = a.logger.Sync()
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErrCh := make(chan error, 1)

	go func() {
		a.logger.Info("starting http server",
			zap.Int("port", a.cfg.Port),
		)
		err := a.httpServer.Start()
		if err != nil && !errors.Is(err, nethttp.ErrServerClosed) {
			serverErrCh <- err
			return
		}
		serverErrCh <- nil
	}()

	select {
	case err := <-serverErrCh:
		if err != nil {
			a.logger.Error("http server stopped with error", zap.Error(err))
		}
		return err

	case <-ctx.Done():
		a.logger.Info("shutdown signal received")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), a.cfg.GracefulShutdownTimeout)
		defer cancel()

		a.logger.Info("starting graceful shutdown")

		if err := a.httpServer.Stop(shutdownCtx); err != nil {
			return fmt.Errorf("failed to stop http server: %w", err)
		}

		a.logger.Info("http server stopped gracefully")

		return nil
	}
}
