package main

import (
	"github.com/ogrock3t/url-shortener-golang/internal/http-server/handlers/redirect"
	"github.com/ogrock3t/url-shortener-golang/internal/http-server/handlers/shorten"
	"github.com/ogrock3t/url-shortener-golang/internal/http-server/middleware/logger"
	"github.com/ogrock3t/url-shortener-golang/internal/service"
	storage "github.com/ogrock3t/url-shortener-golang/internal/storage"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:" + port
	}

	repo := storage.NewInMemoryStorage()
	linkService := service.NewLinkService(repo)

	shortenHandler := shorten.NewHandler(linkService, baseURL)
	redirectHandler := redirect.NewHandler(linkService)

	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", shortenHandler.Handle)
	mux.HandleFunc("/", redirectHandler.Handle)

	loggedMux := logger.MiddlewareLogger(mux)

	log.Printf("server started at :%s", port)

	if err := http.ListenAndServe(":"+port, loggedMux); err != nil {
		log.Fatal(err)
	}
}
