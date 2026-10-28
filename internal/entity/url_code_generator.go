package entity

import (
	"crypto/rand"
	"fmt"
	"io"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomKey() (string, error) {
	return randomKeyFrom(rand.Reader, 7, alphabet)
}

func RandomKey(length int) (string, error) {
	return randomKeyFrom(rand.Reader, length, alphabet)
}

func randomKeyFrom(r io.Reader, length int, chars string) (string, error) {
	if length < 0 {
		return "", fmt.Errorf("length must be non-negative")
	}

	n := len(chars)
	if n < 1 || n > 256 {
		return "", fmt.Errorf("alphabet length must be between 1 and 256")
	}

	// Largest multiple of n that fits in a byte's 256 values.
	limit := 256 - 256%n

	result := make([]byte, length)
	buf := make([]byte, length)

	for i := 0; i < length; {
		if _, err := io.ReadFull(r, buf); err != nil {
			return "", fmt.Errorf("reading random bytes: %w", err)
		}

		for _, b := range buf {
			if int(b) >= limit {
				continue
			}

			result[i] = chars[int(b)%n]
			i++

			if i == length {
				break
			}
		}
	}

	return string(result), nil
}
