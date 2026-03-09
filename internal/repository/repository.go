package repository

import "context"

type LinkRepository interface {
	GetOrCreateID(ctx context.Context, originalURL string) (int64, error)
	GetOriginalURL(ctx context.Context, id int64) (string, error)
}
