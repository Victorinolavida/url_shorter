package httpjson

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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
