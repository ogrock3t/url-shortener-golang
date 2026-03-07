package domain

import (
	"errors"
	"net/url"
	"unicode/utf8"
)

var (
	ErrNotFound    = errors.New("not found")
	ErrInvalidURL  = errors.New("invalid url")
	ErrInvalidCode = errors.New("invalid code")
)

type Link struct {
	ID          int64
	ShortURL    string
	OriginalURL string
}

func ValidateOriginalURL(link string) error {
	if link == "" {
		return ErrInvalidURL
	}

	url, err := url.ParseRequestURI(link)
	if err != nil {
		return ErrInvalidURL
	}

	if url.Scheme != "http" && url.Scheme != "https" {
		return ErrInvalidURL
	}

	if url.Host == "" {
		return ErrInvalidURL
	}

	return nil
}

func ValidateShortURL(link string) error {
	if utf8.RuneCountInString(link) != 10 {
		return ErrInvalidCode
	}

	for _, char := range link {
		if char >= 'a' && char <= 'z' {
			continue
		}

		if char >= 'A' && char <= 'Z' {
			continue
		}

		if char >= '0' && char <= '9' {
			continue
		}

		if char == '_' {
			continue
		}

		return ErrInvalidCode
	}

	return nil
}
