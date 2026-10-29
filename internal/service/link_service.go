package service

import (
	"context"
	"log/slog"

	"url_shorter/internal/entity"
	"url_shorter/internal/repository"
)

type LinkServ struct {
	repo repository.LinkRepository
}

func NewLinkService(repo repository.LinkRepository) LinkService {
	return &LinkServ{repo}
}

func (s *LinkServ) CreateOrGetLinkCode(ctx context.Context, link *entity.Link) (*entity.Link, error) {
	saved, err := s.repo.Save(ctx, link)
	if err != nil {
		return nil, err
	}

	slog.Info("link resolved", slog.String("code", saved.Code))
	return saved, nil
}

func (s *LinkServ) ExpireActiveLinks(ctx context.Context) error {
	count, err := s.repo.ExpireLinks(ctx)
	if err != nil {
		return err
	}
	if count > 0 {
		slog.Info("expired links", slog.Int("count", count))
	}
	return nil
}

func (s *LinkServ) GetActiveLinkByCode(ctx context.Context, code string) (*entity.Link, error) {
	return s.repo.GetActiveByCode(ctx, code)
}
