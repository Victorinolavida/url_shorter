package handler

import (
	"fmt"
	"time"
)

type CreateRequest struct {
	URL       string  `json:"url"`
	ExpiresIn *string `json:"expires_in"`
}

type CreateResponse struct {
	URL       string     `json:"url"`
	ExpiresAt *time.Time `json:"expires_at"`
}

func CalculateExpiration(expiresIn *string) (*time.Time, error) {
	if expiresIn == nil {
		return nil, nil // Never expires
	}

	duration, err := time.ParseDuration(*expiresIn)
	if err != nil {
		return nil, fmt.Errorf("invalid expires_in: %w", err)
	}

	if duration <= 0 {
		return nil, fmt.Errorf("expires_in must be positive")
	}

	expiresAt := time.Now().Add(duration)
	return &expiresAt, nil
}
