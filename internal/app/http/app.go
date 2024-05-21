package httpapp

import (
	"context"
	"fmt"
	"github.com/labstack/gommon/log"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	log        *slog.Logger
	handler    http.Handler
	httpServer http.Server
	port       int
}

// New creates a new HTTP server.
func New(log *slog.Logger, port int, handler http.Handler) *Server {
	return &Server{
		log:     log,
		handler: handler,
		port:    port,
	}
}

// MustRun runs the server and panics if there is an error.
func (s *Server) MustRun() {
	if err := s.Run(); err != nil {
		panic(err)
	}
}

func (s *Server) Run() error {
	const op = "httpapp.Server.Run"
	s.log.With(
		slog.String("op", op),
		slog.Int("port", s.port),
	)
	s.httpServer = http.Server{
		Addr:           ":" + strconv.Itoa(s.port),
		Handler:        s.handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
	}

	if err := s.httpServer.ListenAndServe(); err != nil {

		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("HTTP server is running")
	return nil
}

// Stop stops HTTP server
func (s *Server) Stop(ctx context.Context) {
	const op = "httpapp.Stop"
	s.log.With(slog.String("op", op)).
		Info("stopping http server", slog.Int("port", s.port))
	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		return
	}
}
