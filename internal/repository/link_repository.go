package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"url_shorter/db/sqlc"
	"url_shorter/internal/entity"
	apperrors "url_shorter/internal/errors"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LinkRepo struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

func NewLinkRepository(pool *pgxpool.Pool) LinkRepository {
	return &LinkRepo{
		pool:    pool,
		queries: sqlc.New(pool),
	}
}

func (repo *LinkRepo) Save(ctx context.Context, link *entity.Link) (*entity.Link, error) {
	params := sqlc.SaveLinkParams{
		Code:    link.Code,
		RawUrl:  link.RawUrl,
		UrlHash: link.HashUrl,
	}
	if link.ExpiresAt != nil {
		params.ExpiresAt = pgtype.Timestamptz{Time: link.ExpiresAt.UTC(), Valid: true}
	}

	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	queries := sqlc.New(tx)
	if err := queries.ExpireStaleLinksByHash(ctx, link.HashUrl); err != nil {
		return nil, err
	}

	result, err := queries.SaveLink(ctx, params)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return linkFromRow(result), nil
}

func (repo *LinkRepo) GetByCode(ctx context.Context, code string) (*entity.Link, error) {
	result, err := repo.queries.GetLinkByCode(ctx, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrLinkNotFound
		}
		return nil, err
	}
	return linkFromRow(result), nil
}

func (repo *LinkRepo) GetByHash(ctx context.Context, hash []byte) (*entity.Link, error) {
	result, err := repo.queries.GetLinkByHash(ctx, hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrLinkNotFound
		}
		return nil, err
	}
	return linkFromRow(result), nil
}

func (repo *LinkRepo) GetActiveByHash(ctx context.Context, hash []byte) (*entity.Link, error) {
	result, err := repo.queries.GetActiveLinkByHash(ctx, hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrLinkNotFound
		}
		return nil, err
	}
	return linkFromRow(result), nil
}

func linkFromRow(result sqlc.Link) *entity.Link {
	var expiresAt *time.Time
	if result.ExpiresAt.Valid {
		t := result.ExpiresAt.Time
		expiresAt = &t
	}
	return &entity.Link{
		ID:        int(result.ID),
		Code:      result.Code,
		Status:    entity.LinkStatus(result.Status),
		RawUrl:    result.RawUrl,
		HashUrl:   result.UrlHash,
		CreatedAt: result.CreatedAt.Time,
		ExpiresAt: expiresAt,
	}
}

func (repo *LinkRepo) ExpireLinks(ctx context.Context) (int, error) {
	result, err := repo.queries.ExpireLinks(ctx)
	return int(result), err
}
