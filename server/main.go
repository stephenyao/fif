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

//go:embed dist/*
var staticFiles embed.FS

func getCORSOrigins() []string {
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		// In production, when serving from same origin, we might not need this,
		// but keeping it for flexibility.
		return []string{"http://localhost:5173"}
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

	// API Routes
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

	// Static Files Serving
	publicFS, err := fs.Sub(staticFiles, "dist")
	if err != nil {
		log.Fatal(err)
	}
	fileServer := http.FileServer(http.FS(publicFS))

	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		// If the file exists in the static files, serve it
		f, err := publicFS.Open(strings.TrimPrefix(path, "/"))
		if err == nil {
			f.Close()
			fileServer.ServeHTTP(w, r)
			return
		}

		// Otherwise serve index.html (for React Router)
		index, err := publicFS.Open("index.html")
		if err != nil {
			http.Error(w, "index.html not found", http.StatusNotFound)
			return
		}
		index.Close()
		http.ServeFileFS(w, r, publicFS, "index.html")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
