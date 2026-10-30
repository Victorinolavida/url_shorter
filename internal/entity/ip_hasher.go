package entity

import (
	"crypto/hmac"
	"crypto/sha256"
)

type IPHasher struct{ key []byte }

func NewIPHasher(key []byte) *IPHasher {
	return &IPHasher{key: key}
}

func (h *IPHasher) Hash(ip string) []byte {
	if h == nil || len(h.key) == 0 {
		return nil
	}
	m := hmac.New(sha256.New, h.key)
	m.Write([]byte(ip))
	return m.Sum(nil)
}
