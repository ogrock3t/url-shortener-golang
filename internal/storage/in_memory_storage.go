package memory

import (
	"context"
	"github.com/ogrock3t/url-shortener-golang/internal/domain"
	"sync"
)

type InMemoryStorage struct {
	shortToOriginal map[int64]string
	originalToShort map[string]int64
	count           int64
	mu              sync.RWMutex
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		shortToOriginal: make(map[int64]string),
		originalToShort: make(map[string]int64),
		count:           0,
	}
}

func (s *InMemoryStorage) GetShortURL(ctx context.Context, originalURL string) (int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	code, ok := s.originalToShort[originalURL]
	if !ok {
		return -1, domain.ErrNotFound
	}

	return code, nil
}

func (s *InMemoryStorage) GetOriginalURL(ctx context.Context, id int64) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	url, ok := s.shortToOriginal[id]
	if !ok {
		return "", domain.ErrNotFound
	}

	return url, nil
}

func (s *InMemoryStorage) Save(ctx context.Context, id int64, originalURL string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.shortToOriginal[id] = originalURL
	s.originalToShort[originalURL] = id
	s.count++

	return nil
}

func (s *InMemoryStorage) GetCount(ctx context.Context) int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.count
}
