package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/go-chi/chi/v5"
)

const (
	ctxTokenKey = "idToken" // *auth.Token
)

func main() {
	firebase, err := initFirebaseApp()

	if err != nil {
		log.Fatalf("error initializing firebase app: %v\n", err)
	}

	authClient, err := firebase.Auth(context.Background())

	if err != nil {
		log.Fatalf("error getting auth client: %v\n", err)
	}

	r := chi.NewRouter()

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			jwt := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

			token, err := authClient.VerifyIDToken(r.Context(), jwt)

			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ctxTokenKey, token)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	r.Get("/account", func(w http.ResponseWriter, r *http.Request) {
		token := r.Context().Value(ctxTokenKey).(*auth.Token)
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

	log.Printf("auth client: %v\n", authClient)
	log.Fatal(http.ListenAndServe(":8080", r))
}
