package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/Lovodia/restapi/internal/handlers"
	"github.com/Lovodia/restapi/internal/storage"
	"github.com/Lovodia/restapi/pkg/config"
	loggerSwitch "github.com/Lovodia/restapi/pkg/logger"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		tempLogger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
		tempLogger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	logger := loggerSwitch.NewLogger(cfg.Logger.Level)
	store := storage.NewResultStore()
	e := echo.New()

	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())

	e.POST("/sum", handlers.PostHandler(logger, store))
	e.POST("/multiply", handlers.MultiplyHandler(logger, store))
	e.GET("/results", handlers.GetAllResultsByTokenHandler(logger, store))

	go func() {
		if err := e.Start(":" + cfg.Server.Port); err != nil && err != http.ErrServerClosed {
			logger.Error("shutting down the server due to error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	logger.Info("shutdown signal received", "signal", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logger.Error("error during server shutdown", "error", err)
	} else {
		logger.Info("server shutdown completed gracefully")
	}
}
