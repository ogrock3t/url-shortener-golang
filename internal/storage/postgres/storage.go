package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ogrock3t/url-shortener-golang/internal/domain"
)

type Storage struct {
	pool *pgxpool.Pool
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{
		pool: pool,
	}
}

func (s *Storage) GetOrCreateID(ctx context.Context, originalURL string) (int64, error) {
	const query = `
		INSERT INTO links (original_url)
		VALUES ($1)
		ON CONFLICT (original_url)
		DO UPDATE SET original_url = EXCLUDED.original_url
		RETURNING id
	`

	var id int64

	err := s.pool.QueryRow(ctx, query, originalURL).Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("get or create id: %w", err)
	}

	return id, nil
}

func (s *Storage) GetOriginalURL(ctx context.Context, id int64) (string, error) {
	const query = `
		SELECT original_url
		FROM links
		WHERE id = $1
	`

	var originalURL string

	err := s.pool.QueryRow(ctx, query, id).Scan(&originalURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", domain.ErrNotFound
		}
		return "", fmt.Errorf("get original url: %w", err)
	}

	return originalURL, nil
}
