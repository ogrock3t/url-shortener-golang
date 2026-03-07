package repository

import "context"

type LinkRepository interface {
	GetShortURL(ctx context.Context, originalURL string) (int64, error)
	GetOriginalURL(ctx context.Context, id int64) (string, error)
	Save(ctx context.Context, id int64, originalURL string) error
	GetCount(ctx context.Context) int64
}
