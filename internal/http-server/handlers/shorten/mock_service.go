package shorten

import "context"

type mockLinkService struct {
	createShortURLFn func(ctx context.Context, originalURL string) (string, error)
}

func (m *mockLinkService) CreateShortURL(ctx context.Context, originalURL string) (string, error) {
	return m.createShortURLFn(ctx, originalURL)
}
