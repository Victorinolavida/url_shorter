package service

import (
	"context"
	"errors"
	"log/slog"

	"url_shorter/internal/entity"
	apperrors "url_shorter/internal/errors"
	"url_shorter/internal/repository"
)

type LinkServ struct {
	repo repository.LinkRepository
}

func NewLinkService(repo repository.LinkRepository) LinkService {
	return &LinkServ{repo}
}

func (s *LinkServ) CreateOrGetLinkCode(ctx context.Context, link *entity.Link) (*entity.Link, error) {
	// first check if the url is active
	urlFound, err := s.repo.GetActiveByHash(ctx, link.HashUrl)
	// if there is and erro, check that this error is not because the link does not exist
	if err != nil && !errors.Is(err, apperrors.ErrLinkNotFound) {
		return nil, err
	}
	// this means that we have already create this link so we reuse it if is not expired
	if urlFound != nil {
		if urlFound.IsActive() {
			slog.Info("found active link")
			return urlFound, nil
		}
	}

	// in this point the url does not exist or is already expired , so we can create a new one
	err = s.repo.Save(ctx, link)
	if err != nil {
		return nil, err
	}
	slog.Info("create new link")
	return link, nil
}
