package shorten

import (
	"bytes"
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
		createShortURLFn: func(ctx context.Context, originalURL string) (string, error) {
			if originalURL != "https://google.com" {
				t.Fatalf("Expected originalURL %q, got %q", "https://google.com", originalURL)
			}
			return "AAAAAAAAAB", nil
		},
	}

	handler := NewHandler(service, "http://localhost:8080")

	body, err := json.Marshal(Request{
		URL: "https://google.com",
	})
	if err != nil {
		t.Fatalf("UnExpected marshal error: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	// Act
	handler.Handle(rec, req)

	// Assert
	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var resp Response
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("UnExpected unmarshal error: %v", err)
	}

	Expected := "http://localhost:8080/AAAAAAAAAB"
	if resp.ShortURL != Expected {
		t.Fatalf("Expected short_url %q, got %q", Expected, resp.ShortURL)
	}
}

func TestHandle_MethodNotAllowed(t *testing.T) {
	// Arrange
	service := &mockLinkService{
		createShortURLFn: func(ctx context.Context, originalURL string) (string, error) {
			return "AAAAAAAAAB", nil
		},
	}

	handler := NewHandler(service, "http://localhost:8080")
	req := httptest.NewRequest(http.MethodGet, "/shorten", nil)
	rec := httptest.NewRecorder()

	// Act
	handler.Handle(rec, req)

	// Assert
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("Expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
	}

	var resp ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("UnExpected unmarshal error: %v", err)
	}

	if resp.Error != "Method Not Allowed" {
		t.Fatalf("Expected error %q, got %q", "Method Not Allowed", resp.Error)
	}
}

func TestHandle_InvalidJSON(t *testing.T) {
	// Arrange
	service := &mockLinkService{
		createShortURLFn: func(ctx context.Context, originalURL string) (string, error) {
			return "AAAAAAAAAB", nil
		},
	}

	handler := NewHandler(service, "http://localhost:8080")
	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBufferString("{"))
	rec := httptest.NewRecorder()

	// Act
	handler.Handle(rec, req)

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var resp ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("UnExpected unmarshal error: %v", err)
	}

	if resp.Error == "" {
		t.Fatal("Expected non-empty error")
	}
}

func TestHandle_ServiceError(t *testing.T) {
	// Arrange
	ExpectedErr := errors.New("invalid url")

	service := &mockLinkService{
		createShortURLFn: func(ctx context.Context, originalURL string) (string, error) {
			return "", ExpectedErr
		},
	}

	handler := NewHandler(service, "http://localhost:8080")

	body, err := json.Marshal(Request{
		URL: "invalid",
	})
	if err != nil {
		t.Fatalf("UnExpected marshal error: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	// Act
	handler.Handle(rec, req)

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var resp ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("UnExpected unmarshal error: %v", err)
	}

	if resp.Error != ExpectedErr.Error() {
		t.Fatalf("Expected error %q, got %q", ExpectedErr.Error(), resp.Error)
	}
}
