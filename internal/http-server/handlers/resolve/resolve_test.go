package resolve

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
			if shortURL != "AAAAAAAAAA" {
				t.Fatalf("Expected shortURL %q, got %q", "AAAAAAAAAA", shortURL)
			}

			return "https://google.com", nil
		},
	}

	handler := NewHandler(service)
	req := httptest.NewRequest(http.MethodGet, "/resolve/AAAAAAAAAA", nil)
	rec := httptest.NewRecorder()

	// Act
	handler.Handle(rec, req)

	// Assert
	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var resp Response
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Unexpected unmarshal error: %v", err)
	}

	if resp.OriginalURL != "https://google.com" {
		t.Fatalf("Expected original_url %q, got %q", "https://google.com", resp.OriginalURL)
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
	req := httptest.NewRequest(http.MethodPost, "/resolve/AAAAAAAAAA", nil)
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
	req := httptest.NewRequest(http.MethodGet, "/resolve/", nil)
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
	expectedErr := errors.New("Not found")

	service := &mockLinkService{
		findOriginalURLFn: func(ctx context.Context, shortURL string) (string, error) {
			return "", expectedErr
		},
	}

	handler := NewHandler(service)
	req := httptest.NewRequest(http.MethodGet, "/resolve/AAAAAAAAAA", nil)
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

	if resp.Error != expectedErr.Error() {
		t.Fatalf("Expected error %q, got %q", expectedErr.Error(), resp.Error)
	}
}
