package memory

import (
	"context"
	"github.com/ogrock3t/url-shortener-golang/internal/domain"
	"sync"
)

type InMemoryStorage struct {
	storage map[int64]string
	count   int64
	mu      sync.RWMutex
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		storage: make(map[int64]string),
		count:   0,
	}
}

func (s *InMemoryStorage) GetShortURL(ctx context.Context, originalURL string) (int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for id, url := range s.storage {
		if originalURL == url {
			return id, nil
		}
	}

	return -1, domain.ErrNotFound
}

func (s *InMemoryStorage) GetOriginalURL(ctx context.Context, id int64) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	url, ok := s.storage[id]
	if !ok {
		return "", domain.ErrNotFound
	}

	return url, nil
}

func (s *InMemoryStorage) Save(ctx context.Context, id int64, originalURL string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.storage[id] = originalURL
	s.count++

	return nil
}

func (s *InMemoryStorage) GetCount(ctx context.Context) int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.count
}
