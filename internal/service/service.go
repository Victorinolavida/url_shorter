package service

import (
	"context"

	"url_shorter/internal/entity"
)

type LinkService interface {
	CreateOrGetLinkCode(context.Context, *entity.Link) (*entity.Link, error)
	ExpireActiveLinks(context.Context) error
	GetActiveLinkByCode(context.Context, string) (*entity.Link, error)
	RecordLinkClickEvent(context.Context, *entity.ClickEvent) error
}
