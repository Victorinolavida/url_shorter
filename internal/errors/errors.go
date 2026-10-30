package apperrors

import "errors"

var (
	ErrLinkNotFound      = errors.New("link not found")
	ErrLinkExpired       = errors.New("link expired")
	ErrLinkDisabled      = errors.New("link disabled")
	ErrLinkBlocked       = errors.New("link blocked")
	ErrInvalidURLSchema  = errors.New("invalid schema for url")
	ErrInvalidHostname   = errors.New("invalid hostname")
	ErrInvalidPort       = errors.New("invalid port")
	ErrNoUserInfo        = errors.New("no user info found")
	ErrUrlToBig          = errors.New("url to big")
	ErrNoPortFound       = errors.New("no port info found")
	ErrUrlEmpty          = errors.New("url must be not empty")
	ErrInvalidExpiration = errors.New("invalid expiration time")
	ErrInvalidLinkID     = errors.New("invalid link id")
	ErrUserAgentTooLong  = errors.New("user agent is too long")
	ErrInvalidStatusCode = errors.New("invalid status code")
)
