package main

import (
	"context"
	"encoding/json"
	"fif/middleware"
	"log"
	"net/http"
	"os"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type HoldingDTO struct {
	Name     string  `json:"name"`
	Symbol   string  `json:"symbol"`
	Quantity float64 `json:"quantity"`
	Currency string  `json:"currency"`
	Cost     float64 `json:"cost"`
}

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

	r.Use(middleware.AuthMiddleware(authClient))

	r.Get("/account", func(w http.ResponseWriter, r *http.Request) {
		token := r.Context().Value(middleware.CtxTokenKey{}).(*auth.Token)
		email, _ := token.Claims["email"].(string)
		name, _ := token.Claims["name"].(string)

		resp := map[string]string{
			"email": email,
			"name":  name,
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	r.Get("/holdings", func(w http.ResponseWriter, r *http.Request) {
		resp := []HoldingDTO{
			{
				Name:     "Block",
				Symbol:   "XYZ",
				Quantity: 10,
				Currency: "USD",
				Cost:     123.45,
			},
			{
				Name:     "Google",
				Symbol:   "GOOG",
				Quantity: 30,
				Currency: "USD",
				Cost:     223.45,
			},
			{
				Name:     "Apple",
				Symbol:   "APPL",
				Quantity: 40,
				Currency: "USD",
				Cost:     139.45,
			},
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
