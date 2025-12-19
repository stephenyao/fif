package middleware

import (
	"context"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
)

// CtxTokenKey is the context key for storing the auth token
type CtxTokenKey struct{}

// AuthVerifier is an interface for verifying authentication tokens
type AuthVerifier interface {
	VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error)
}

// firebaseAuthClient wraps the Firebase auth.Client to implement AuthVerifier
type firebaseAuthClient struct {
	client *auth.Client
}

func (f *firebaseAuthClient) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	return f.client.VerifyIDToken(ctx, idToken)
}

func AuthMiddleware(client *auth.Client) func(handler http.Handler) http.Handler {
	return authMiddlewareWithVerifier(&firebaseAuthClient{client: client})
}

func authMiddlewareWithVerifier(verifier AuthVerifier) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			jwt := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := verifier.VerifyIDToken(r.Context(), jwt)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), CtxTokenKey{}, token)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
