package redirect

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

type LinkService interface {
	FindOriginalURL(ctx context.Context, originalURL string) (string, error)
}

type Handler struct {
	service LinkService
}

func NewHandler(service LinkService) *Handler {
	return &Handler{
		service: service,
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, &ErrorResponse{Error: "Method Not Allowed"})
		return
	}

	code := strings.TrimPrefix(r.URL.Path, "/")
	if code == "" {
		http.NotFound(w, r)
		return
	}

	originalURL, err := h.service.FindOriginalURL(r.Context(), code)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, &ErrorResponse{Error: err.Error()})
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}

func writeJSON(w http.ResponseWriter, status int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}
