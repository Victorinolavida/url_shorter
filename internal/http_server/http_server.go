package httpserver

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"url_shorter/config"
	"url_shorter/internal/httpjson"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

var defaultPort = 8080

type ServerInterface interface {
	Start() error
	Shutdown(context.Context) error
	Mux() *chi.Mux
}

var _ ServerInterface = (*httpServer)(nil)

type httpServer struct {
	srv    *http.Server
	router *chi.Mux
}

func NewHttpServer(conf config.Server) ServerInterface {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if reqID := middleware.GetReqID(r.Context()); reqID != "" {
				w.Header().Set(middleware.RequestIDHeader, reqID)
			}
			next.ServeHTTP(w, r)
		})
	})
	r.Use(middleware.ClientIPFromHeader("X-Real-IP"))
	r.Use(middleware.Recoverer)

	// The browser session is cookie based, so the cors handler has to allow
	// credentials. Browsers reject "*" together with credentials, therefore the
	// allowed origins have to be listed explicitly.

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of the major browsers
	}))

	srv := &http.Server{
		Addr:        formatServerAddress(conf.Host, conf.Port),
		Handler:     r,
		ReadTimeout: time.Duration(10 * time.Second),
		IdleTimeout: time.Duration(100 * time.Second),
	}

	httpSvc := &httpServer{
		srv:    srv,
		router: r,
	}

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		httpjson.WriteNotFound(w)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		httpjson.WriteNotMethodNotAllowed(w)
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return httpSvc
}

// Mux returns the underlying chi router
func (h *httpServer) Mux() *chi.Mux {
	return h.router
}

func (h *httpServer) Start() error {
	slog.Info("HTTP server starting", slog.String("addr", h.srv.Addr))
	return h.srv.ListenAndServe()
}

func (h *httpServer) Shutdown(ctx context.Context) error {
	slog.Info("HTTP server shutting down...")
	return h.srv.Shutdown(ctx)
}

func formatServerAddress(host string, port int) string {
	if port == 0 {
		port = defaultPort
	}
	if host == "" {
		host = "localhost"
	}
	return host + ":" + strconv.Itoa(port)
}
