package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/brpaz/echozap"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func main() {
	// ==== Config ====
	cfg, err := LoadConfig()
	expect(err == nil, "Load config:", err)

	// ==== Logger ====
	var logger *zap.Logger
	switch cfg.Logger {
	case "dev":
		logger, err = zap.NewDevelopment()
	case "prod":
		logger, err = zap.NewProduction()
	default:
		err = errors.New("logger `" + cfg.Logger + "` is unknown")
	}
	expect(err == nil, "Logger:", err)

	defer func() {
		_ = logger.Sync()
	}()
	sugaredLogger := logger.Sugar()

	// ==== HTTP Server ====
	e := echo.New()
	e.HideBanner = true
	e.Use(echozap.ZapLogger(logger))

	// ==== Setup handlers ====
	handler := NewHandler(sugaredLogger, e)
	handler.RegisterRoutes()

	// ==== Start server ====
	go func() {
		sugaredLogger.Info("Starting server...")

		if err = e.Start(":" + strconv.Itoa(cfg.Port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
			sugaredLogger.Fatalf("Shutting down the server: %s", err)
		}
	}()

	// ==== Graceful shutdown ====
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	sugaredLogger.Info("Stopping server...")
	if err = e.Shutdown(ctx); err != nil {
		sugaredLogger.Fatalf("Shutdown: %s", err)
	}
}

func expect(ok bool, args ...any) {
	if !ok {
		fmt.Println(append([]any{"FATAL:"}, args...)...)
		os.Exit(1)
	}
}
