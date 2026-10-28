package entity

import (
	"errors"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"
)

type User struct {
	ID        int
	Code      string
	CreatedAt time.Time
	URL       string
}

var (
	ErrInvalidSchema   = errors.New("invalid schema for url")
	ErrInvalidHostname = errors.New("invalid hostname")
	ErrInvalidPort     = errors.New("invalid port")
	ErrNoUserInfo      = errors.New("no user info found")
	ErrUrlToBig        = errors.New("url to big")
	ErrNoPortFound     = errors.New("no port info found")

	ErrUrlEmpty    = errors.New("url must be not empty")
	urlLimitLength = 2048
	ValidSchema    = []string{"http", "https"}
	portRange      = [2]int{1, 65535}
)

func NewUser(
	rawurl string,
) (*User, error) {
	trimedUrl := strings.TrimSpace(rawurl)
	if trimedUrl == "" {
		return nil, ErrUrlEmpty
	}
	if len(trimedUrl) > urlLimitLength {
		return nil, ErrUrlToBig
	}
	url, err := url.ParseRequestURI(trimedUrl)
	if err != nil {
		return nil, err
	}
	if !slices.Contains(ValidSchema, url.Scheme) {
		return nil, ErrInvalidSchema
	}
	if url.Hostname() == "" {
		return nil, ErrInvalidHostname
	}
	if url.Port() != "" {
		port, err := strconv.Atoi(url.Port())
		if err != nil {
			return nil, err
		}
		if port < portRange[0] || port > portRange[1] {
			return nil, ErrInvalidPort
		}
	}
	if url.User != nil {
		return nil, ErrNoUserInfo
	}

	code, err := GenerateRandomKey()
	if err != nil {
		return nil, err
	}

	return &User{
		URL:  rawurl,
		Code: code,
	}, nil
}
