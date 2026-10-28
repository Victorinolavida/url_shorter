package entity

import (
	"crypto/sha256"
	"net/url"
	"slices"
	"strconv"
	"strings"

	apperrors "url_shorter/internal/errors"
)

type Url struct {
	Raw  string
	Hash []byte
}

var (
	urlLimitLength = 2048
	ValidSchema    = []string{"http", "https"}
	portRange      = [2]int{1, 65535}
)

func NewURL(raw string) (*Url, error) {
	trimedUrl := strings.TrimSpace(raw)
	if trimedUrl == "" {
		return nil, apperrors.ErrUrlEmpty
	}
	if len(trimedUrl) > urlLimitLength {
		return nil, apperrors.ErrUrlToBig
	}
	url, err := url.ParseRequestURI(trimedUrl)
	if err != nil {
		return nil, err
	}
	if !slices.Contains(ValidSchema, url.Scheme) {
		return nil, apperrors.ErrInvalidURLSchema
	}
	if url.Hostname() == "" {
		return nil, apperrors.ErrInvalidHostname
	}
	if url.Port() != "" {
		port, err := strconv.Atoi(url.Port())
		if err != nil {
			return nil, err
		}
		if port < portRange[0] || port > portRange[1] {
			return nil, apperrors.ErrInvalidPort
		}
	}
	if url.User != nil {
		return nil, apperrors.ErrNoUserInfo
	}
	hash := sha256.Sum256([]byte(raw))
	return &Url{
		Raw:  raw,
		Hash: hash[:],
	}, nil
}
