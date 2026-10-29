package entity

import (
	"time"

	apperrors "url_shorter/internal/errors"
)

type Link struct {
	ID        int
	Code      string
	Status    LinkStatus
	ExpiresAt *time.Time
	CreatedAt time.Time
	RawUrl    string
	HashUrl   []byte
}

// IsActive return true if the status is active
func (l *Link) IsActive() bool {
	return l.Status == ActiveStatus
}

func (l *Link) IsExpired() bool {
	return l.ExpiresAt != nil && l.ExpiresAt.Before(time.Now())
}

func NewLink(raw string, expired *time.Time) (*Link, error) {
	url, err := NewURL(raw)
	if err != nil {
		return nil, err
	}

	code, err := GenerateRandomKey()
	if err != nil {
		return nil, err
	}
	if expired != nil && expired.Before(time.Now()) {
		return nil, apperrors.ErrInvalidExpiration
	}
	return &Link{
		Code:      code,
		ExpiresAt: expired,
		HashUrl:   url.Hash,
		RawUrl:    url.Raw,
	}, nil
}
