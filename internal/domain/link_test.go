package domain

import (
	"errors"
	"testing"
)

func TestValidateOriginalURL_OK(t *testing.T) {
	// Arrange
	validURLs := []string{
		"http://google.com",
		"https://youtu.be/8-daMyeA8RE?si=4WYhrb3H81d1e4_C",
		"https://github.com/ogrock3t",
	}

	// Act && Assert
	for _, url := range validURLs {
		err := ValidateOriginalURL(url)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}
	}
}

func TestValidateOriginalURL_Invalid(t *testing.T) {
	// Arrange
	invalidURL := []string{
		"google.com",
		"",
		"invalid",
	}

	// Act && Assert
	for _, url := range invalidURL {
		err := ValidateOriginalURL(url)
		if !errors.Is(err, ErrInvalidURL) {
			t.Fatalf("Expected ErrInvalidURL for %q, got %v", url, err)
		}
	}
}

func TestValidateShortURL_OK(t *testing.T) {
	// Arrange
	validCodes := []string{
		"1234567890",
		"AAAAAAAAAA",
		"abcd___234",
	}

	// Act && Assert
	for _, code := range validCodes {
		err := ValidateShortURL(code)
		if err != nil {
			t.Fatalf("Expected no error for %q, got %v", code, err)
		}
	}
}

func TestValidateOriginalURL_InvalidURL(t *testing.T) {
	// Arrange
	invalidCodes := []string{
		"",
		"abc",
		"abcd_-_234",
	}

	// Act && Assert
	for _, code := range invalidCodes {
		err := ValidateShortURL(code)
		if !errors.Is(err, ErrInvalidCode) {
			t.Fatalf("Expected ErrInvalidCode for %q, got %v", code, err)
		}
	}
}
