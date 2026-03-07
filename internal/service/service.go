package service

import (
	"context"
	"github.com/ogrock3t/url-shortener-golang/internal/domain"
	"github.com/ogrock3t/url-shortener-golang/internal/repository"
)

type LinkService struct {
	repository repository.LinkRepository
}

func NewLinkService(repository repository.LinkRepository) *LinkService {
	return &LinkService{repository: repository}
}

func (service *LinkService) CreateShortURL(ctx context.Context, originalURL string) (string, error) {
	if err := domain.ValidateOriginalURL(originalURL); err != nil {
		return "", err
	}

	if id, err := service.repository.GetShortURL(ctx, originalURL); err == nil {
		return GenerateShortURL(id), nil
	} else if err != nil && err != domain.ErrNotFound {
		return "", err
	}

	id := service.repository.GetCount(ctx)
	shortURL := GenerateShortURL(id)

	if err := service.repository.Save(ctx, id, originalURL); err != nil {
		return "", err
	}

	return shortURL, nil
}

func (service *LinkService) FindOriginalURL(ctx context.Context, shortURL string) (string, error) {
	if err := domain.ValidateShortURL(shortURL); err != nil {
		return "", err
	}

	id := GetIDFromShortURL(shortURL)

	return service.repository.GetOriginalURL(ctx, id)
}

func GenerateShortURL(id int64) string {
	alphabet := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789" + "_")

	result := make([]byte, 10)
	index := 9

	for id > 0 {
		result[index] = alphabet[id%63]
		id /= 63
		index--
	}

	left := index + 1
	right := 9
	for left < right {
		result[left], result[right] = result[right], result[left]
		left++
		right--
	}

	for index >= 0 {
		result[index] = 'A'
		index--
	}

	return string(result)
}

func GetIDFromShortURL(shortURL string) int64 {
	alphabet := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789" + "_")

	id := int64(0)

	for i := 0; i < len(shortURL); i++ {
		char := shortURL[i]

		var value int64

		for j := 0; j < len(alphabet); j++ {
			if alphabet[j] == char {
				value = int64(j)
				break
			}
		}

		id = id*63 + value
	}

	return id
}
