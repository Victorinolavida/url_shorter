package cache

import (
	"context"
	"errors"

	"url_shorter/internal/entity"
)

var ErrCacheMiss = errors.New("cache: key is missing")

type Cache interface {
	Get(context.Context, string) (*entity.Link, error)
	Set(context.Context, *entity.Link) error
	Delete(context.Context, string) error
}
