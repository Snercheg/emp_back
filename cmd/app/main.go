package main

import (
	"EMP_Back/internal/app"
	"EMP_Back/internal/config"
	"fmt"
	"github.com/joho/godotenv"
	"golang.org/x/net/context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// @title EMP Back
// @version 1.0.0
// @description API Server for EMP backend application

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	// TODO: Change config delete console
	cfg := config.MustLoadConfig()
	fmt.Println(cfg)
	// TODO: Customize logger
	if err := godotenv.Load(); err != nil {
		slog.Error("Error loading.env file", slog.String("error", err.Error()))
	}

	log := setupLogger(cfg.Env)
	log.Info("Starting application",
		slog.String("env", cfg.Env),
	)

	log.Debug("Debug message")
	log.Error("Error message")
	log.Warn("Warning message")

	application := app.New(log, cfg.Server.Port, os.Getenv("DB_URL"), cfg.TokenTTL)
	go application.HTTPServer.MustRun()

	// TODO: setup database

	// TODO: middleware

	// TODO: setup handlers

	// TODO: setup routes

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	sign := <-stop
	log.Info("Shutting down...", slog.String("signal", sign.String()))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	application.HTTPServer.Stop(ctx)

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
