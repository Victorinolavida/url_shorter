package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"url_shorter/internal/entity"
	"url_shorter/internal/httpjson"
	"url_shorter/internal/service"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service service.LinkService
}

func NewHandler(service service.LinkService) *Handler {
	return &Handler{service}
}

func (h *Handler) HandleCreateCode(w http.ResponseWriter, r *http.Request) {
	payload, err := httpjson.DecodeRequest[CreateRequest](w, r)
	if err != nil {
		log.Printf("decode request: %v", err)

		switch {
		case errors.Is(err, httpjson.ErrBodyTooLarge):
			http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
		case errors.Is(err, httpjson.ErrEmptyBody),
			errors.Is(err, httpjson.ErrInvalidJSON):
			http.Error(w, "Invalid request body", http.StatusBadRequest)
		default:
			http.Error(w, "Bad request", http.StatusBadRequest)
		}
		return
	}
	expiredAt, err := CalculateExpiration(payload.ExpiresIn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	link, err := entity.NewLink(payload.URL, expiredAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.service.CreateOrGetLinkCode(r.Context(), link)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := CreateResponse{
		URL:       "http://localhost:5100/" + result.Code,
		ExpiresAt: result.ExpiresAt,
	}

	buf, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(buf)
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/url-shorten", h.HandleCreateCode)
	// r.Get("/{code}", h.Redirect)
}
