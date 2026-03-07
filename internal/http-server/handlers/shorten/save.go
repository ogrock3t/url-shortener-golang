package shorten

import (
	"context"
	"encoding/json"
	"net/http"
)

type LinkService interface {
	CreateShortURL(ctx context.Context, originalURL string) (string, error)
}

type Handler struct {
	service LinkService
	baseURL string
}

func NewHandler(service LinkService, baseURL string) *Handler {
	return &Handler{
		service: service,
		baseURL: baseURL,
	}
}

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	ShortURL string `json:"short_url"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, &ErrorResponse{Error: "Method Not Allowed"})
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, &ErrorResponse{Error: err.Error()})
		return
	}

	code, err := h.service.CreateShortURL(r.Context(), req.URL)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, &ErrorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, &Response{
		ShortURL: h.baseURL + "/" + code,
	})
}

func writeJSON(w http.ResponseWriter, status int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}
