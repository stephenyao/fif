package handlers

import (
	"context"
	"encoding/json"
	"fif/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"firebase.google.com/go/v4/auth"
)

func TestAccountHandler_Success(t *testing.T) {
	// Create a mock token with claims
	token := &auth.Token{
		UID: "test-user-123",
		Claims: map[string]interface{}{
			"email": "test@example.com",
			"name":  "Test User",
		},
	}

	// Create a request with the token in context
	req := httptest.NewRequest(http.MethodGet, "/account", nil)
	ctx := context.WithValue(req.Context(), middleware.CtxTokenKey{}, token)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	// Execute the handler
	AccountHandler(w, req)

	// Assert status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Assert Content-Type header
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	// Parse and assert response body
	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp["email"] != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", resp["email"])
	}

	if resp["name"] != "Test User" {
		t.Errorf("Expected name 'Test User', got %s", resp["name"])
	}
}

func TestAccountHandler_MissingToken(t *testing.T) {
	// Create a request without a token in context
	req := httptest.NewRequest(http.MethodGet, "/account", nil)
	w := httptest.NewRecorder()

	// Execute the handler
	AccountHandler(w, req)

	// Assert status code
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	// Assert response body
	expectedBody := "unauthorized\n"
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, w.Body.String())
	}
}

func TestAccountHandler_NilToken(t *testing.T) {
	// Create a request with nil token in context
	req := httptest.NewRequest(http.MethodGet, "/account", nil)
	ctx := context.WithValue(req.Context(), middleware.CtxTokenKey{}, nil)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	// Execute the handler
	AccountHandler(w, req)

	// Assert status code
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	// Assert response body
	expectedBody := "unauthorized\n"
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, w.Body.String())
	}
}

func TestAccountHandler_WrongTypeInContext(t *testing.T) {
	// Create a request with wrong type in context (string instead of *auth.Token)
	req := httptest.NewRequest(http.MethodGet, "/account", nil)
	ctx := context.WithValue(req.Context(), middleware.CtxTokenKey{}, "not-a-token")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	// Execute the handler
	AccountHandler(w, req)

	// Assert status code
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAccountHandler_MissingEmailClaim(t *testing.T) {
	// Create a token without email claim
	token := &auth.Token{
		UID: "test-user-123",
		Claims: map[string]interface{}{
			"name": "Test User",
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/account", nil)
	ctx := context.WithValue(req.Context(), middleware.CtxTokenKey{}, token)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	// Execute the handler
	AccountHandler(w, req)

	// Should still succeed with empty email
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp["email"] != "" {
		t.Errorf("Expected empty email, got %s", resp["email"])
	}

	if resp["name"] != "Test User" {
		t.Errorf("Expected name 'Test User', got %s", resp["name"])
	}
}

func TestAccountHandler_MissingNameClaim(t *testing.T) {
	// Create a token without name claim
	token := &auth.Token{
		UID: "test-user-123",
		Claims: map[string]interface{}{
			"email": "test@example.com",
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/account", nil)
	ctx := context.WithValue(req.Context(), middleware.CtxTokenKey{}, token)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	// Execute the handler
	AccountHandler(w, req)

	// Should still succeed with empty name
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp["email"] != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", resp["email"])
	}

	if resp["name"] != "" {
		t.Errorf("Expected empty name, got %s", resp["name"])
	}
}

func TestAccountHandler_NonStringClaims(t *testing.T) {
	// Create a token with non-string claim values
	token := &auth.Token{
		UID: "test-user-123",
		Claims: map[string]interface{}{
			"email": 12345, // Wrong type
			"name":  true,  // Wrong type
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/account", nil)
	ctx := context.WithValue(req.Context(), middleware.CtxTokenKey{}, token)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	// Execute the handler
	AccountHandler(w, req)

	// Should still succeed but with empty values (type assertion fails)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Both should be empty since type assertion fails
	if resp["email"] != "" {
		t.Errorf("Expected empty email for non-string claim, got %s", resp["email"])
	}

	if resp["name"] != "" {
		t.Errorf("Expected empty name for non-string claim, got %s", resp["name"])
	}
}
