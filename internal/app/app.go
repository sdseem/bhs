// Package app configures and runs application.
package app

import (
	"bhs/internal/usecase/assets"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"bhs/config"
	v1 "bhs/internal/controller/http/v1"
	"bhs/internal/usecase/auth"
	"bhs/internal/usecase/repo"
	"bhs/pkg/httpserver"
	"bhs/pkg/logger"
	"bhs/pkg/postgres"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Use case
	authUseCase := auth.NewAuth(repo.NewUserRepo(pg), cfg.App.TokenSecret)
	assetsUseCase := assets.NewAssets(repo.NewAssetsRepo(pg))

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, l, authUseCase, assetsUseCase)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
