package entity

import (
	"slices"
	"strings"
	"time"

	apperrors "url_shorter/internal/errors"
)

const maxUserAgentSize = 512

var validStatusCode = []int{301, 303, 302, 410}

type ClickEvent struct {
	ID         int
	LinkID     int
	ClickedAt  time.Time
	UserAgent  string
	Browser    string
	Os         string
	DeviceType string
	Country    string
	IPHash     []byte
	IsBot      bool
	StatusCode int
}

func NewClickEvent(
	linkID int,
	userAgent, browser, osReq, deviceType,
	country, ip string, isBot bool, statusCode int,
	hasher *IPHasher,
) (*ClickEvent, error) {
	if linkID <= 0 {
		return nil, apperrors.ErrInvalidLinkID
	}
	if !slices.Contains(validStatusCode, statusCode) {
		return nil, apperrors.ErrInvalidStatusCode
	}
	var ipHash []byte
	if hasher != nil && ip != "" {
		ipHash = hasher.Hash(ip)
	}
	return &ClickEvent{
		LinkID:     linkID,
		UserAgent:  truncateUA(userAgent),
		Browser:    browser,
		Os:         osReq,
		IsBot:      isBot,
		DeviceType: deviceType,
		Country:    country,
		IPHash:     ipHash,
		StatusCode: statusCode,
	}, nil
}

func truncateUA(s string) string {
	if len(s) <= maxUserAgentSize {
		return s
	}
	// cut on a byte boundary, then drop any broken trailing UTF-8
	return strings.ToValidUTF8(s[:maxUserAgentSize], "")
}
