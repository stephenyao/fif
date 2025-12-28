package handlers

import (
	"database/sql"
	"encoding/json"
	"fif/middleware"
	"log"
	"net/http"

	"firebase.google.com/go/v4/auth"
)

// HoldingDTO represents a financial holding
type HoldingDTO struct {
	Name     string  `json:"name"`
	Symbol   string  `json:"symbol"`
	Quantity float64 `json:"quantity"`
	Currency string  `json:"currency"`
	Cost     float64 `json:"cost"`
}

// MakeHoldingsHandler creates a handler that fetches holdings from the database
func MakeHoldingsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the auth token from context
		token, ok := r.Context().Value(middleware.CtxTokenKey{}).(*auth.Token)
		if !ok || token == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		userID := token.UID

		// Query holdings for this user
		rows, err := db.Query(`
			SELECT name, symbol, quantity, currency, cost
			FROM holdings
			WHERE user_id = $1
			ORDER BY created_at DESC
		`, userID)
		if err != nil {
			log.Printf("Error querying holdings: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		holdings := []HoldingDTO{}
		for rows.Next() {
			var h HoldingDTO
			if err := rows.Scan(&h.Name, &h.Symbol, &h.Quantity, &h.Currency, &h.Cost); err != nil {
				log.Printf("Error scanning holding: %v", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			holdings = append(holdings, h)
		}

		if err := rows.Err(); err != nil {
			log.Printf("Error iterating holdings: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(holdings); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
