package repository

import (
	"context"
	"database/sql"
	"errors"

	"url_shorter/db/sqlc"
	"url_shorter/internal/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UrlCodeRepo struct {
	queries *sqlc.Queries
}

var ErrURLNotFound = errors.New("url not found")

func NewUrlCodeRepository(pool *pgxpool.Pool) UrlCodeRepository {
	return &UrlCodeRepo{
		queries: sqlc.New(pool),
	}
}

func (repo *UrlCodeRepo) Save(ctx context.Context, user *entity.User) error {
	result, err := repo.queries.SaveUrlCode(ctx, sqlc.SaveUrlCodeParams{
		Code: user.Code,
		Url:  user.URL,
	})
	if err != nil {
		return err
	}
	user.ID = int(result.ID)
	user.CreatedAt = result.CreatedAt.Time
	return nil
}

func (repo *UrlCodeRepo) GetByCode(ctx context.Context, code string) (*entity.User, error) {
	result, err := repo.queries.GetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrURLNotFound
		}
		return nil, err
	}
	return &entity.User{
		ID:        int(result.ID),
		Code:      result.Code,
		URL:       result.Url,
		CreatedAt: result.CreatedAt.Time,
	}, nil
}

func (repo *UrlCodeRepo) GetByUrl(ctx context.Context, code string) (*entity.User, error) {
	result, err := repo.queries.GetByURL(ctx, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrURLNotFound
		}
		return nil, err
	}
	return &entity.User{
		ID:        int(result.ID),
		Code:      result.Code,
		URL:       result.Url,
		CreatedAt: result.CreatedAt.Time,
	}, nil
}
