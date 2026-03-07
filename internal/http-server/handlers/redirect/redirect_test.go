package redirect

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandle_OK(t *testing.T) {
	// Arrange
	service := &mockLinkService{
		findOriginalURLFn: func(ctx context.Context, shortURL string) (string, error) {
			if shortURL != "AAAAAAAAAB" {
				t.Fatalf("Expected shortURL %q, got %q", "AAAAAAAAAB", shortURL)
			}
			return "https://google.com", nil
		},
	}

	handler := NewHandler(service)
	req := httptest.NewRequest(http.MethodGet, "/AAAAAAAAAB", nil)
	rec := httptest.NewRecorder()

	// Act
	handler.Handle(rec, req)

	// Assert
	if rec.Code != http.StatusFound {
		t.Fatalf("Expected status %d, got %d", http.StatusFound, rec.Code)
	}

	location := rec.Header().Get("Location")
	if location != "https://google.com" {
		t.Fatalf("Expected location %q, got %q", "https://google.com", location)
	}
}

func TestHandle_MethodNotAllowed(t *testing.T) {
	// Arrange
	service := &mockLinkService{
		findOriginalURLFn: func(ctx context.Context, shortURL string) (string, error) {
			return "https://google.com", nil
		},
	}

	handler := NewHandler(service)
	req := httptest.NewRequest(http.MethodPost, "/AAAAAAAAAB", nil)
	rec := httptest.NewRecorder()

	// Act
	handler.Handle(rec, req)

	// Assert
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("Expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
	}

	var resp ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Unexpected unmarshal error: %v", err)
	}

	if resp.Error != "Method Not Allowed" {
		t.Fatalf("Expected error %q, got %q", "Method Not Allowed", resp.Error)
	}
}

func TestHandle_EmptyCode(t *testing.T) {
	// Arrange
	service := &mockLinkService{
		findOriginalURLFn: func(ctx context.Context, shortURL string) (string, error) {
			return "https://google.com", nil
		},
	}

	handler := NewHandler(service)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	// Act
	handler.Handle(rec, req)

	// Assert
	if rec.Code != http.StatusNotFound {
		t.Fatalf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestHandle_ServiceError(t *testing.T) {
	// Arrange
	ExpectedErr := errors.New("not found")

	service := &mockLinkService{
		findOriginalURLFn: func(ctx context.Context, shortURL string) (string, error) {
			return "", ExpectedErr
		},
	}

	handler := NewHandler(service)
	req := httptest.NewRequest(http.MethodGet, "/AAAAAAAAAB", nil)
	rec := httptest.NewRecorder()

	// Act
	handler.Handle(rec, req)

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var resp ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Unexpected unmarshal error: %v", err)
	}

	if resp.Error != ExpectedErr.Error() {
		t.Fatalf("Expected error %q, got %q", ExpectedErr.Error(), resp.Error)
	}
}
