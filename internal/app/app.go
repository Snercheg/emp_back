package app

import (
	httpapp "EMP_Back/internal/app/http"
	"EMP_Back/internal/http-server/handler"
	"EMP_Back/internal/service"
	"EMP_Back/internal/storage"
	"log/slog"
	"time"
)

type App struct {
	HTTPServer *httpapp.Server
}

func New(
	log *slog.Logger,
	port int,
	dbURL string,
	tokenTTL time.Duration,
) *App {
	// TODO init storage
	db, err := storage.New(storage.DBConfig{
		DBUrl: dbURL,
	})
	if err != nil {
		log.Error("Error in init storage", err)
	}
	repos := storage.NewRepository(db)
	services := service.NewService(repos, log)
	handlers := handler.NewHandler(services)

	// TODO init app service
	// handlers := new(handler.Handler)

	httpApp := httpapp.New(log, port, handlers.InitRoutes())
	return &App{
		HTTPServer: httpApp,
	}
}
