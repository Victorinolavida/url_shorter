package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) CreateURL(
	w http.ResponseWriter, r *http.Request,
) {
}

func (h *Handler) Reditect(
	w http.ResponseWriter, r *http.Request,
) {
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewMux()
	r.Post("/api/shorten", h.CreateURL)
	r.Get("/{code}", h.Reditect)
	return r
}
