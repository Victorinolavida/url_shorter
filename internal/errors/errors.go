package apperrors

import "errors"

var (
	ErrLinkNotFound      = errors.New("link not found")
	ErrInvalidURLSchema  = errors.New("invalid schema for url")
	ErrInvalidHostname   = errors.New("invalid hostname")
	ErrInvalidPort       = errors.New("invalid port")
	ErrNoUserInfo        = errors.New("no user info found")
	ErrUrlToBig          = errors.New("url to big")
	ErrNoPortFound       = errors.New("no port info found")
	ErrUrlEmpty          = errors.New("url must be not empty")
	ErrInvalidExpiration = errors.New("invalid expiration time")
)
