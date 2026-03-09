package in_memory

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

func (s *InMemoryStorage) GetOrCreateID(_ context.Context, originalURL string) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, url := range s.storage {
		if url == originalURL {
			return id, nil
		}
	}

	id := s.count
	s.storage[id] = originalURL
	s.count++

	return id, nil
}

func (s *InMemoryStorage) GetOriginalURL(_ context.Context, id int64) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	url, ok := s.storage[id]
	if !ok {
		return "", domain.ErrNotFound
	}

	return url, nil
}
