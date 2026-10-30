package service

import (
	"context"
	"errors"
	"log/slog"

	"url_shorter/cache"
	"url_shorter/internal/entity"
	"url_shorter/internal/repository"
)

type LinkServ struct {
	repo  repository.LinkRepository
	cache cache.Cache
}

func NewLinkService(repo repository.LinkRepository, cache cache.Cache) LinkService {
	return &LinkServ{repo, cache}
}

func (s *LinkServ) CreateOrGetLinkCode(ctx context.Context, link *entity.Link) (*entity.Link, error) {
	saved, err := s.repo.Save(ctx, link)
	if err != nil {
		return nil, err
	}

	slog.Info("link resolved", slog.String("code", saved.Code))
	if err := s.cache.Set(ctx, saved); err != nil {
		slog.Warn("unable to save in cache", slog.Any("error", err))
	}
	return saved, nil
}

func (s *LinkServ) ExpireActiveLinks(ctx context.Context) error {
	codes, err := s.repo.ExpireLinks(ctx)
	if err != nil {
		return err
	}
	for _, code := range codes {
		if err := s.cache.Delete(ctx, code); err != nil {
			slog.Warn("unable to delete expired link from cache", slog.String("code", code), slog.Any("error", err))
		}
	}
	if len(codes) > 0 {
		slog.Info("expired links", slog.Int("count", len(codes)))
	}
	return nil
}

func (s *LinkServ) GetActiveLinkByCode(ctx context.Context, code string) (*entity.Link, error) {
	cached, err := s.cache.Get(ctx, code)
	switch {
	case err == nil:
		if cached == nil {
			break
		}
		if err := cached.ValidateActive(); err != nil {
			slog.Warn("cached link is not active", slog.String("code", code), slog.Any("reason", err))
			if err := s.cache.Delete(ctx, code); err != nil {
				slog.Warn("unable to delete stale cache", slog.String("code", code), slog.Any("error", err))
			}
			break
		}
		return cached, nil
	case errors.Is(err, cache.ErrCacheMiss):
		// Cache miss; fall back to the database.
	default:
		slog.Warn("error getting from cache", slog.String("code", code), slog.Any("error", err))
	}

	link, err := s.repo.GetActiveByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if err := s.cache.Set(ctx, link); err != nil {
		slog.Warn("unable to save in cache", slog.Any("error", err))
	}
	return link, nil
}

func (s *LinkServ) RecordLinkClickEvent(ctx context.Context, event *entity.ClickEvent) error {
	return s.repo.SaveClickEvent(ctx, event)
}
