package main

import (
	"EMP_Back/internal/config"
	"fmt"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// TODO: initialize config object
	cfg := config.MustLoadConfig()
	fmt.Println(cfg)
	// TODO: Initialize logger

	log := setupLogger(cfg.Env)
	log.Info("Starting gRPC server",
		slog.String("env", cfg.Env),
	)

	log.Debug("Debug message")
	log.Error("Error message")
	log.Warn("Warning message")

	// TODO: Initialize application

	// TODO: Run gRPC server application

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
