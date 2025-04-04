package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

var (
	mu    sync.Mutex
	store = make(map[string]string)
)

func generateID() string {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return "errorid"
	}
	return hex.EncodeToString(b)
}

func handlePost(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.Error(w, "Invalid URL for POST", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		http.Error(w, "Empty request body", http.StatusBadRequest)
		return
	}
	origURL := strings.TrimSpace(string(body))
	if origURL == "" {
		http.Error(w, "Empty URL provided", http.StatusBadRequest)
		return
	}

	id := generateID()

	mu.Lock()
	store[id] = origURL
	mu.Unlock()

	shortURL := fmt.Sprintf("http://localhost:8080/%s", id)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated) // 201 Created
	_, _ = w.Write([]byte(shortURL))
}

func handleGet(w http.ResponseWriter, r *http.Request) {

	id := strings.TrimPrefix(r.URL.Path, "/")
	if id == "" {
		http.Error(w, "Missing short URL id", http.StatusBadRequest)
		return
	}
	mu.Lock()
	origURL, ok := store[id]
	mu.Unlock()
	if !ok {
		http.Error(w, "Short URL not found", http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", origURL)
	w.WriteHeader(http.StatusTemporaryRedirect) // 307 Temporary Redirect
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handlePost(w, r)
	case http.MethodGet:
		handleGet(w, r)
	default:
		http.Error(w, "Unsupported method", http.StatusBadRequest)
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("Server is listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
