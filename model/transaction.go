package model

import (
	"math"
	"time"
)

const expectedDateFormat = "2006-01-02"

type Transaction struct {
	ID              int     `json:"id"`
	Description     string  `json:"description"`
	Amount          float64 `json:"amount"`
	TransactionDate string  `json:"transaction_date"`
}

// Validate checks the Transaction fields for validity.
func (t *Transaction) Validate() string {
	if len(t.Description) > 50 {
		return "Description must be 50 characters or fewer"
	}

	if t.Amount <= 0 {
		return "Amount must be greater than 0"
	}

	t.Amount = roundToCents(t.Amount)

	if _, err := time.Parse(expectedDateFormat, t.TransactionDate); err != nil {
		return "Transaction date must be in YYYY-MM-DD format"
	}

	return ""
}

// roundToCents rounds a float64 value to the nearest cent.
func roundToCents(value float64) float64 {
	return math.Round(value*100) / 100
}
