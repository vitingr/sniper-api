package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"main.go/config"
	"main.go/internal/di"
	logger "main.go/pkg/logger"
	router "main.go/internal/infra/http/router"
)

func main() {
	godotenv.Load()

	cfg := config.Load()

	logger := logger.New(cfg.LogLevel)

	container := di.NewContainer()

	routerWrapper := router.NewRouter(container)
	ginHandler := routerWrapper.Setup()

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      ginHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("Starting server",
			"port", cfg.Port,
			"env", cfg.AppEnv,
		)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", "error", err)
	}

	logger.Info("Server exited gracefully")
}
