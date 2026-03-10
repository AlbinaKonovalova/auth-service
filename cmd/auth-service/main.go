package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/AlbinaKonovalova/auth-service/internal/config"
	"github.com/AlbinaKonovalova/auth-service/internal/infrastructure/app/authservice"
)

func main() {
	configPath := flag.String("config", "config/config.yaml", "path to config file")
	flag.Parse()

	logger := setupLogger("info")
	logger.Info("starting auth-service")

	cfg, err := config.Load(*configPath)
	if err != nil {
		logger.Error("failed to load config", slog.Any("error", err))
		os.Exit(1)
	}

	logger = setupLogger(cfg.Log.Level)

	app, err := authservice.New(cfg, logger)
	if err != nil {
		logger.Error("failed to initialize app", slog.Any("error", err))
		os.Exit(1)
	}

	go func() {
		if err := app.Run(); err != nil {
			logger.Error("server error", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", slog.Any("error", err))
	}

	logger.Info("server stopped")
}

func setupLogger(level string) *slog.Logger {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
}
