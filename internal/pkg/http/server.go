package http

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"time"

	"url_shorter/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	slogchi "github.com/samber/slog-chi"
)

var (
	defaultPort = 8080
	defaultHost = "localhost"
)

type Server struct {
	serv   *http.Server
	router *chi.Mux
}

func buildAddres(host string, port int) string {
	if host == "" {
		host = defaultHost
	}
	if port <= 0 {
		port = defaultPort
	}

	portStr := strconv.Itoa(port)
	return net.JoinHostPort(host, portStr)
}

func NewServer(conf config.Server) *Server {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(slogchi.New(slog.Default()))
	r.Use(middleware.CleanPath)
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.Recoverer)

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	serv := &http.Server{
		Handler:           r,
		Addr:              buildAddres(conf.Host, conf.Port),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       75 * time.Second,
	}

	return &Server{
		serv:   serv,
		router: r,
	}
}

func (s *Server) Start() error {
	slog.Info("listening server on", slog.String("addr", s.serv.Addr))
	return s.serv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	slog.Info("shutdown server....")
	return s.serv.Shutdown(ctx)
}

// Mux returns the underlying chi router
func (h *Server) Mux() *chi.Mux {
	return h.router
}
