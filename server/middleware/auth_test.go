package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"firebase.google.com/go/v4/auth"
)

// mockAuthVerifier is a mock implementation of the AuthVerifier interface for testing
type mockAuthVerifier struct {
	verifyFunc func(ctx context.Context, idToken string) (*auth.Token, error)
}

func (m *mockAuthVerifier) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	if m.verifyFunc != nil {
		return m.verifyFunc(ctx, idToken)
	}
	return nil, errors.New("not implemented")
}

// mockHandler is a simple handler that writes a success message
func mockHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})
}

func TestAuthMiddleware_MissingAuthorizationHeader(t *testing.T) {
	// Create a mock auth verifier (won't be called in this test)
	mockVerifier := &mockAuthVerifier{}

	// Create middleware
	middleware := AuthMiddlewareWithVerifier(mockVerifier)
	handler := middleware(mockHandler())

	// Create a request without Authorization header
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	// Execute
	handler.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	expectedBody := "unauthorized\n"
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, w.Body.String())
	}
}

func TestAuthMiddleware_InvalidAuthorizationHeaderFormat(t *testing.T) {
	testCases := []struct {
		name        string
		authHeader  string
		description string
	}{
		{
			name:        "NoBearer",
			authHeader:  "InvalidToken",
			description: "Authorization header without 'Bearer ' prefix",
		},
		{
			name:        "WrongPrefix",
			authHeader:  "Basic sometoken",
			description: "Authorization header with wrong prefix",
		},
		{
			name:        "EmptyHeader",
			authHeader:  "",
			description: "Empty authorization header",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockVerifier := &mockAuthVerifier{}
			middleware := AuthMiddlewareWithVerifier(mockVerifier)
			handler := middleware(mockHandler())

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != http.StatusUnauthorized {
				t.Errorf("Expected status %d for %s, got %d", http.StatusUnauthorized, tc.description, w.Code)
			}
		})
	}
}

func TestAuthMiddleware_EmptyTokenAfterBearer(t *testing.T) {
	mockVerifier := &mockAuthVerifier{
		verifyFunc: func(ctx context.Context, idToken string) (*auth.Token, error) {
			// This will be called with an empty string
			if idToken == "" {
				return nil, errors.New("empty token")
			}
			return &auth.Token{UID: "test"}, nil
		},
	}

	middleware := AuthMiddlewareWithVerifier(mockVerifier)
	handler := middleware(mockHandler())

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer ")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Should fail with unauthorized
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddleware_TokenVerificationFailure(t *testing.T) {
	// Create a mock verifier that returns an error
	mockVerifier := &mockAuthVerifier{
		verifyFunc: func(ctx context.Context, idToken string) (*auth.Token, error) {
			return nil, errors.New("invalid token")
		},
	}

	middleware := AuthMiddlewareWithVerifier(mockVerifier)
	handler := middleware(mockHandler())

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	expectedBody := "unauthorized\n"
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, w.Body.String())
	}
}

func TestAuthMiddleware_SuccessfulAuthentication(t *testing.T) {
	// Create a valid token for testing
	expectedToken := &auth.Token{
		UID: "test-user-123",
		Claims: map[string]interface{}{
			"email": "test@example.com",
		},
	}

	// Create a mock verifier that returns a valid token
	mockVerifier := &mockAuthVerifier{
		verifyFunc: func(ctx context.Context, idToken string) (*auth.Token, error) {
			if idToken == "valid-token" {
				return expectedToken, nil
			}
			return nil, errors.New("invalid token")
		},
	}

	middleware := AuthMiddlewareWithVerifier(mockVerifier)

	// Create a handler that checks if the token was added to context
	var contextToken *auth.Token
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from context
		token := r.Context().Value(ctxTokenKeyType{})
		if token != nil {
			contextToken = token.(*auth.Token)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	handler := middleware(testHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Assert HTTP response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	expectedBody := "success"
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, w.Body.String())
	}

	// Assert token was added to context
	if contextToken == nil {
		t.Fatal("Expected token to be added to context, got nil")
	}

	if contextToken.UID != expectedToken.UID {
		t.Errorf("Expected UID %s, got %s", expectedToken.UID, contextToken.UID)
	}
}

func TestAuthMiddleware_TokenExtraction(t *testing.T) {
	var capturedToken string

	// Create a mock verifier that captures the token passed to it
	mockVerifier := &mockAuthVerifier{
		verifyFunc: func(ctx context.Context, idToken string) (*auth.Token, error) {
			capturedToken = idToken
			return &auth.Token{UID: "test-user"}, nil
		},
	}

	middleware := AuthMiddlewareWithVerifier(mockVerifier)
	handler := middleware(mockHandler())

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer my-test-token")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Verify the "Bearer " prefix was stripped correctly
	expectedToken := "my-test-token"
	if capturedToken != expectedToken {
		t.Errorf("Expected token %q to be passed to VerifyIDToken, got %q", expectedToken, capturedToken)
	}
}

func TestAuthMiddleware_ContextPropagation(t *testing.T) {
	// Create a mock verifier
	mockVerifier := &mockAuthVerifier{
		verifyFunc: func(ctx context.Context, idToken string) (*auth.Token, error) {
			// Verify that the original request context is passed through
			if ctx == nil {
				t.Error("Expected context to be non-nil")
			}
			return &auth.Token{UID: "test-user"}, nil
		},
	}

	middleware := AuthMiddlewareWithVerifier(mockVerifier)

	// Create a handler that checks the context
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify that context is not nil
		if r.Context() == nil {
			t.Error("Expected request context to be non-nil")
		}

		// Verify token is in context
		token := r.Context().Value(ctxTokenKeyType{})
		if token == nil {
			t.Error("Expected token in context, got nil")
		}

		w.WriteHeader(http.StatusOK)
	})

	handler := middleware(testHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestAuthMiddleware_BearerWithExtraSpaces(t *testing.T) {
	var capturedToken string

	mockVerifier := &mockAuthVerifier{
		verifyFunc: func(ctx context.Context, idToken string) (*auth.Token, error) {
			capturedToken = idToken
			return &auth.Token{UID: "test-user"}, nil
		},
	}

	middleware := AuthMiddlewareWithVerifier(mockVerifier)
	handler := middleware(mockHandler())

	// Test with token that has leading/trailing spaces after "Bearer "
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer  token-with-spaces ")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// The middleware strips "Bearer " but doesn't trim the token
	// So extra spaces will be included
	expectedToken := " token-with-spaces "
	if capturedToken != expectedToken {
		t.Errorf("Expected token %q, got %q", expectedToken, capturedToken)
	}
}

func TestAuthMiddleware_MultipleRequests(t *testing.T) {
	callCount := 0
	mockVerifier := &mockAuthVerifier{
		verifyFunc: func(ctx context.Context, idToken string) (*auth.Token, error) {
			callCount++
			return &auth.Token{UID: "test-user"}, nil
		},
	}

	middleware := AuthMiddlewareWithVerifier(mockVerifier)
	handler := middleware(mockHandler())

	// Make multiple requests
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer token")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Request %d: Expected status %d, got %d", i+1, http.StatusOK, w.Code)
		}
	}

	// Verify the verifier was called for each request
	if callCount != 3 {
		t.Errorf("Expected verifier to be called 3 times, got %d", callCount)
	}
}
