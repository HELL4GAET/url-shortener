package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strings"
	"time"
)

type URLService interface {
	Shorten(ctx context.Context, originalURL string) (string, error)
	Resolve(ctx context.Context, alias string) (string, error)
	Stats(ctx context.Context, alias string) (int, time.Time, string, error)
}

type URLHandler struct {
	service URLService
}

func NewURLHandler(service URLService) *URLHandler {
	return &URLHandler{service: service}
}

func (h *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	alias := chi.URLParam(r, "code")
	if alias == "" {
		http.Error(w, "Missing short URL", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	originalURL, err := h.service.Resolve(ctx, alias)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
		originalURL = "https://" + originalURL
	}

	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}

func (h *URLHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	originalURL := strings.TrimSpace(string(body))
	if originalURL == "" {
		http.Error(w, "empty URL provided", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	alias, err := h.service.Shorten(ctx, originalURL)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	shortURL := fmt.Sprintf("http://localhost:8080/%s", alias)
	w.Header().Set("Content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(shortURL))
}

func (h *URLHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	alias := chi.URLParam(r, "code")
	if alias == "" {
		http.Error(w, "Missing short URL", http.StatusBadRequest)
		return
	}
	clicks, createdAt, originalURL, err := h.service.Stats(ctx, alias)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(
		map[string]interface{}{
			"original_url": originalURL,
			"clicks":       clicks,
			"created_at":   createdAt.Format(time.RFC3339),
		})
}
