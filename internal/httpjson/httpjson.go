package httpjson

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	apperrors "url_shorter/internal/errors"
)

const maxRequestBodySize int64 = 8 * 1024 // 8 KiB

var (
	ErrBodyTooLarge = errors.New("request body too large")
	ErrEmptyBody    = errors.New("request body is empty")
	ErrInvalidJSON  = errors.New("invalid JSON request body")
)

func DecodeRequest[T any](w http.ResponseWriter, r *http.Request) (*T, error) {
	if r == nil || r.Body == nil {
		return nil, ErrEmptyBody
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var payload T
	if err := decoder.Decode(&payload); err != nil {
		if errors.Is(err, io.EOF) {
			return nil, ErrEmptyBody
		}

		_, ok := errors.AsType[*http.MaxBytesError](err)
		if ok {
			return nil, ErrBodyTooLarge
		}

		return nil, fmt.Errorf("%w: %w", ErrInvalidJSON, err)
	}

	var extra json.RawMessage
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		_, ok := errors.AsType[*http.MaxBytesError](err)
		if ok {
			return nil, ErrBodyTooLarge
		}

		if err == nil {
			return nil, fmt.Errorf("%w: multiple JSON values", ErrInvalidJSON)
		}

		return nil, fmt.Errorf("%w: %w", ErrInvalidJSON, err)
	}

	return &payload, nil
}

func Envelop(key string, data any) map[string]any {
	return map[string]any{key: data}
}

func WriteJsonResponse(w http.ResponseWriter, data any, status int) error {
	d, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(d)
	return nil
}

func WriteErrResponse(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrBodyTooLarge):
		WriteErrMsg(w, http.StatusRequestEntityTooLarge, "request body too large")

	case errors.Is(err, ErrEmptyBody),
		errors.Is(err, ErrInvalidJSON):
		WriteErrMsg(w, http.StatusBadRequest, "invalid request body")

	case errors.Is(err, apperrors.ErrLinkNotFound),
		errors.Is(err, apperrors.ErrInvalidURLSchema),
		errors.Is(err, apperrors.ErrInvalidHostname),
		errors.Is(err, apperrors.ErrInvalidPort),
		errors.Is(err, apperrors.ErrNoUserInfo),
		errors.Is(err, apperrors.ErrUrlToBig),
		errors.Is(err, apperrors.ErrNoPortFound),
		errors.Is(err, apperrors.ErrUrlEmpty),
		errors.Is(err, apperrors.ErrInvalidExpiration):
		WriteErrMsg(w, http.StatusBadRequest, err.Error())
	default:
		slog.Error("send error to user",
			slog.Any("error", err),
		)
		WriteErrMsg(w, http.StatusInternalServerError, "something went wrong")
	}
}

func WriteErrMsg(w http.ResponseWriter, status int, msg string) {
	envelop := Envelop("error", msg)
	WriteJsonResponse(w, envelop, status)
}

func WriteNotFound(w http.ResponseWriter) {
	WriteErrMsg(w, http.StatusNotFound, "resource not found")
}

func WriteNotMethodNotAllowed(w http.ResponseWriter) {
	WriteErrMsg(w, http.StatusMethodNotAllowed, "method not allowed")
}
