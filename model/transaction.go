package model

import (
	"time"

	"github.com/mvfavila/transactions/config"
	"github.com/mvfavila/transactions/util"
)

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

	t.Amount = util.RoundToCents(t.Amount)

	if _, err := time.Parse(config.AppConfig.ExpectedDateFormat, t.TransactionDate); err != nil {
		return "Transaction date must be in YYYY-MM-DD format"
	}

	return ""
}
