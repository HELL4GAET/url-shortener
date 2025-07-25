package main

import (
	h "URL-shortener/internal/delivery/http"
	"URL-shortener/internal/repository/postgres"
	"URL-shortener/internal/usecase"
	"URL-shortener/pkg/db"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
)

func main() {
	postgresConnectionString := os.Getenv("DATABASE_URL")

	storage, err := db.New(postgresConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer storage.Close()

	urlRepo := postgres.NewURLRepository(storage.DB())
	service := usecase.NewURLService(urlRepo)
	handler := h.NewURLHandler(service)

	r := chi.NewRouter()
	r.Post("/api/v1/shorten", handler.Shorten)
	r.Get("/api/v1/stats/{code}", handler.GetStats)
	r.Get("/{code}", handler.Redirect)

	log.Println("Server is listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
