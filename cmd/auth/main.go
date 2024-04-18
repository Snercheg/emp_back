package main

import (
	"EMP_Back/internal/app"
	"EMP_Back/internal/config"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// TODO: Change config delete console
	cfg := config.MustLoadConfig()
	fmt.Println(cfg)

	// TODO: Customize logger

	log := setupLogger(cfg.Env)
	log.Info("Starting application",
		slog.String("env", cfg.Env),
	)

	log.Debug("Debug message")
	log.Error("Error message")
	log.Warn("Warning message")

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)
	go application.GRPCSrv.MustRun()

	// TODO: Run gRPC server application

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	sign := <-stop
	log.Info("Shutting down...", slog.String("signal", sign.String()))
	application.GRPCSrv.Stop()
	log.Info("Shut down")
	os.Exit(0)
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
