package redirect

import "context"

type mockLinkService struct {
	findOriginalURLFn func(ctx context.Context, shortURL string) (string, error)
}

func (m *mockLinkService) FindOriginalURL(ctx context.Context, shortURL string) (string, error) {
	return m.findOriginalURLFn(ctx, shortURL)
}
