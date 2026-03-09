package main

import (
	"context"
	"flag"
	"github.com/joho/godotenv"
	"github.com/ogrock3t/url-shortener-golang/internal/http-server/handlers/resolve"
	"log"
	"net/http"
	"os"

	"github.com/ogrock3t/url-shortener-golang/internal/http-server/handlers/redirect"
	"github.com/ogrock3t/url-shortener-golang/internal/http-server/handlers/shorten"
	"github.com/ogrock3t/url-shortener-golang/internal/http-server/middleware/logger"
	"github.com/ogrock3t/url-shortener-golang/internal/repository"
	"github.com/ogrock3t/url-shortener-golang/internal/service"
	inmemory "github.com/ogrock3t/url-shortener-golang/internal/storage/in_memory"
	postgres "github.com/ogrock3t/url-shortener-golang/internal/storage/postgres"
)

const (
	defaultBaseURL = "http://localhost:"
	defaultPort    = "8080"
)

type Config struct {
	Port        string
	BaseURL     string
	StorageType string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found")
	}

	cfg := loadConfig()
	ctx := context.Background()

	repo := newRepository(ctx, cfg)
	handler := newHTTPHandler(cfg, repo)

	log.Printf("Server started at: %s", cfg.Port)

	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		log.Fatal(err)
	}
}

func loadConfig() Config {
	storageType := flag.String("storage", "in-memory", "storage type: in-memory or postgres")
	flag.Parse()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = defaultBaseURL + port
	}

	return Config{
		Port:        port,
		BaseURL:     baseURL,
		StorageType: *storageType,
	}
}

func newRepository(ctx context.Context, cfg Config) repository.LinkRepository {
	switch cfg.StorageType {
	case "in-memory":
		log.Println("using in-memory storage")
		return inmemory.NewInMemoryStorage()

	case "postgres":
		databaseURL := os.Getenv("DATABASE_URL")
		if databaseURL == "" {
			log.Fatal("DATABASE_URL is required for postgres storage")
		}

		pool, err := postgres.NewPool(ctx, databaseURL)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("using postgres storage")

		return postgres.NewStorage(pool)

	default:
		log.Fatalf("unknown storage type: %s", cfg.StorageType)
		return nil
	}
}

func newHTTPHandler(cfg Config, repo repository.LinkRepository) http.Handler {
	linkService := service.NewLinkService(repo)

	shortenHandler := shorten.NewHandler(linkService, cfg.BaseURL)
	resolveHandler := resolve.NewHandler(linkService)
	redirectHandler := redirect.NewHandler(linkService)

	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", shortenHandler.Handle)
	mux.HandleFunc("/resolve/", resolveHandler.Handle)
	mux.HandleFunc("/", redirectHandler.Handle)

	return logger.MiddlewareLogger(mux)
}
