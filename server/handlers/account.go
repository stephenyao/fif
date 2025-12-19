package handlers

import (
	"encoding/json"
	"fif/middleware"
	"net/http"

	"firebase.google.com/go/v4/auth"
)

// AccountHandler handles the /account endpoint
func AccountHandler(w http.ResponseWriter, r *http.Request) {
	token, ok := r.Context().Value(middleware.CtxTokenKey{}).(*auth.Token)
	if !ok || token == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

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
}
