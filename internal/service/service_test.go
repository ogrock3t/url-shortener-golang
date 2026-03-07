package service

import (
	"context"
	storage "github.com/ogrock3t/url-shortener-golang/internal/storage"
	"testing"
	"unicode/utf8"
)

func TestCreateShortURL_OK(t *testing.T) {
	// Arrange
	repo := storage.NewInMemoryStorage()
	service := NewLinkService(repo)

	url := "http://google.com"

	// Act
	code, _ := service.CreateShortURL(context.Background(), url)

	// Assert
	if utf8.RuneCountInString(code) != 10 {
		t.Fatalf("Expected 10 length, got %d", utf8.RuneCountInString(code))
	}
}

func TestCreateShortURL_Idempotent(t *testing.T) {
	// Arrange
	repo := storage.NewInMemoryStorage()
	service := NewLinkService(repo)

	url := "http://google.com"

	// Act
	code1, _ := service.CreateShortURL(context.Background(), url)
	code2, _ := service.CreateShortURL(context.Background(), url)

	// Assert
	if code1 != code2 {
		t.Fatalf("Expected code %s, got %s", code1, code2)
	}
}

func TestFindOriginalURL_OK(t *testing.T) {
	// Arrange
	repo := storage.NewInMemoryStorage()
	service := NewLinkService(repo)

	url := "http://google.com"

	// Act
	code, _ := service.CreateShortURL(context.Background(), url)
	result, _ := service.FindOriginalURL(context.Background(), code)

	// Assert
	if result != url {
		t.Fatalf("Expected %s, got %s", url, result)
	}
}

func TestCreateShortURL_InvalidURL(t *testing.T) {
	// Arrange
	repo := storage.NewInMemoryStorage()
	service := NewLinkService(repo)

	invalidURL := "invalidURL"

	// Act
	_, err := service.CreateShortURL(context.Background(), invalidURL)

	// Assert
	if err == nil {
		t.Fatalf("Expected error, got none")
	}
}
