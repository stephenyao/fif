package handlers

import (
	"encoding/json"
	"net/http"
)

// HoldingDTO represents a financial holding
type HoldingDTO struct {
	Name     string  `json:"name"`
	Symbol   string  `json:"symbol"`
	Quantity float64 `json:"quantity"`
	Currency string  `json:"currency"`
	Cost     float64 `json:"cost"`
}

// HoldingsHandler handles the /holdings endpoint
func HoldingsHandler(w http.ResponseWriter, r *http.Request) {
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
}
