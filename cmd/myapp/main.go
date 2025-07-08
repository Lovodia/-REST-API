// Package main API сервер для вычислений

// Документация Swagger.

// 	Schemest: http https
// 	Host: localhost:8080
// 	BasePath: /
// 	Version: 1.0.0
// 	Contact: Your Name

// 	Consumes:
// 	-application/json

// 	Produces:
// 	-application/json

// Swagger:meta
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/Lovodia/restapi/docs"
	_ "github.com/Lovodia/restapi/docs"
	"github.com/Lovodia/restapi/internal/handlers"
	"github.com/Lovodia/restapi/internal/storage"
	"github.com/Lovodia/restapi/pkg/config"
	loggerswitch "github.com/Lovodia/restapi/pkg/logger"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/swaggo/files"
)

// @title CalculatorAPI
// @version 1.0
// @description API для вычисления суммы и произведения
// @host localhost:8080
// @BasePath /
func main() {
	docs.SwaggerInfo.Schemes = []string{"http"}
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger := loggerswitch.NewLogger(cfg.Logger.Level)
	store := storage.NewResultStore()
	e := echo.New()

	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())

	e.POST("/sum", handlers.PostHandler(logger, store))
	e.POST("/multiply", handlers.MultiplyHandler(logger, store))
	e.GET("/results", handlers.GetAllResultsByTokenHandler(logger, store))
	e.GET("/swagger/*", echoSwagger.WrapHandler)

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
