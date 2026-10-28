package repository

import (
	"context"

	"url_shorter/internal/entity"
)

type UrlCodeRepository interface {
	Save(context.Context, *entity.User) error
	GetByCode(context.Context, string) (*entity.User, error)
	GetByUrl(context.Context, string) (*entity.User, error)
}
