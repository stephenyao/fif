package main

import (
	"context"
	"fif/handlers"
	"fif/middleware"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

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

	// Protected routes (authentication required)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(authClient))

		r.Get("/account", handlers.AccountHandler)

		r.Get("/holdings", handlers.HoldingsHandler)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
