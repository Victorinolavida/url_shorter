package handler

import (
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
		httpjson.WriteErrResponse(w, err)
		return
	}
	expiredAt, err := CalculateExpiration(payload.ExpiresIn)
	if err != nil {
		httpjson.WriteErrResponse(w, err)
		return
	}
	link, err := entity.NewLink(payload.URL, expiredAt)
	if err != nil {
		httpjson.WriteErrResponse(w, err)
		return
	}

	result, err := h.service.CreateOrGetLinkCode(r.Context(), link)
	if err != nil {
		httpjson.WriteErrResponse(w, err)
		return
	}

	response := CreateResponse{
		URL:       "http://localhost:5100/" + result.Code,
		ExpiresAt: result.ExpiresAt,
	}
	httpjson.WriteJsonResponse(w, response, http.StatusOK)
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/url-shorten", h.HandleCreateCode)
	// r.Get("/{code}", h.Redirect)
}
