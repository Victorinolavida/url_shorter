package repository

import (
	"context"

	"url_shorter/internal/entity"
)

type LinkRepository interface {
	Save(context.Context, *entity.Link) (*entity.Link, error)
	GetByCode(context.Context, string) (*entity.Link, error)
	GetByHash(context.Context, []byte) (*entity.Link, error)
	GetActiveByHash(context.Context, []byte) (*entity.Link, error)
	ExpireLinks(context.Context) (int, error)
}
