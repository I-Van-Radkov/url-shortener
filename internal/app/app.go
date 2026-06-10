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
	"github.com/gin-gonic/gin"
)

type App struct {
	httpServer *httpserver.Server
	cfg        *config.Config
	db         *db.Database
}

func NewApp(cfg *config.Config) (*App, error) {
	a := &App{
		cfg: cfg,
	}

	gen := usecase.NewCodeGenerator()

	var uc *usecase.Usecase

	switch a.cfg.StorageType {
	case "memory":
		memRepo := memory.NewRepo()
		uc = usecase.NewUsecase(memRepo, gen, a.cfg.MaxAttemptsToGen)
	case "postgres":
		database, err := db.NewPostgres(cfg.PostgresConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}

		a.db = database

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
			a.db.Close()
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErrCh := make(chan error, 1)

	go func() {
		err := a.httpServer.Start()
		if err != nil && !errors.Is(err, nethttp.ErrServerClosed) {
			serverErrCh <- err
			return
		}
		serverErrCh <- nil
	}()

	select {
	case err := <-serverErrCh:
		return err

	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), a.cfg.GracefulShutdownTimeout)
		defer cancel()

		if err := a.httpServer.Stop(shutdownCtx); err != nil {
			return fmt.Errorf("failed to stop http server: %w", err)
		}

		return nil
	}
}
