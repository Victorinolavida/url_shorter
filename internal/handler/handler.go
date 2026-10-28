package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"url_shorter/internal/entity"
	"url_shorter/internal/repository"

	"github.com/go-chi/chi"
)

type Handler struct {
	repo repository.UrlCodeRepository
}

func NewHandler(repo repository.UrlCodeRepository) *Handler {
	return &Handler{repo}
}

func (h *Handler) CreateCode(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	payload := new(CreateRequest)
	encoder := json.NewDecoder(r.Body)

	if err := encoder.Decode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	urlCode, err := entity.NewUser(payload.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	urlFound, err := h.repo.GetByUrl(r.Context(), payload.URL)
	if err != nil && !errors.Is(err, repository.ErrURLNotFound) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if urlFound != nil {
		slog.Info("url already created")
		response := CreateResponse{URL: "http://localhost:5100/" + urlFound.Code}
		buf, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(buf)

		return
	}

	if err := h.repo.Save(r.Context(), urlCode); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := CreateResponse{URL: "http://localhost:5100/" + urlCode.Code}

	buf, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(buf)
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		http.Error(w, "code is required", http.StatusBadRequest)
		return
	}
	urlCode, err := h.repo.GetByCode(r.Context(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, urlCode.URL, http.StatusSeeOther)
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Post("/url-shorten", h.CreateCode)
	r.Get("/{code}", h.Redirect)
	return r
}
