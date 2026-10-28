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

// IsActive return true if the status is active and the expiration is less than time.Now
// is the expiration is null, is enough that the status is active
func (l *Link) IsActive() bool {
	if l.Status != ActiveStatus {
		return false
	}
	if l.ExpiresAt != nil && l.ExpiresAt.Before(time.Now()) {
		return false
	}
	return true
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
