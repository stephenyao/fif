package main

import (
	"context"
	"embed"
	"fif/handlers"
	"fif/middleware"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

//go:embed all:webdist/*
var webdist embed.FS

func getCORSOrigins() []string {
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		log.Fatal("ALLOWED_ORIGINS environment variable is required")
	}

	origins := strings.Split(allowedOrigins, ",")
	// Trim whitespace from each origin
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}
	return origins
}

func main() {
	firebase, err := initFirebaseApp()

	if err != nil {
		log.Fatalf("error initializing firebase app: %v\n", err)
	}

	authClient, err := firebase.Auth(context.Background())

	if err != nil {
		log.Fatalf("error getting auth client: %v\n", err)
	}

	origins := getCORSOrigins()

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"ok"}`))
		})

		// Protected routes (authentication required)
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(authClient))

			r.Get("/account", handlers.AccountHandler)

			r.Get("/holdings", handlers.HoldingsHandler)
		})
	})

	// Static files and SPA fallback
	distFS, err := fs.Sub(webdist, "webdist")
	if err != nil {
		log.Fatal(err)
	}

	r.Get("/*", handlers.SPAHandler(distFS))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
