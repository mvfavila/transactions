package model

import (
	"testing"
)

func TestTransactionValidate(t *testing.T) {
	tests := []struct {
		name           string
		transaction    Transaction
		expectedResult string
	}{
		{
			name: "Description longer than 50 characters",
			transaction: Transaction{
				Description:     "This is exactly 51 chars long. Oooooops, too long!!",
				Amount:          1.00,
				TransactionDate: "2020-01-01",
			},
			expectedResult: "Description must be 50 characters or fewer",
		},
		{
			name: "Amount less than 0",
			transaction: Transaction{
				Description:     "Test",
				Amount:          -0.01,
				TransactionDate: "2020-01-01",
			},
			expectedResult: "Amount must be greater than 0",
		},
		{
			name: "Amount equals to 0",
			transaction: Transaction{
				Description:     "Test",
				Amount:          0.00,
				TransactionDate: "2020-01-01",
			},
			expectedResult: "Amount must be greater than 0",
		},
		{
			name: "Invalid date format",
			transaction: Transaction{
				Description:     "Test",
				Amount:          1.00,
				TransactionDate: "01-01-2020",
			},
			expectedResult: "Transaction date must be in YYYY-MM-DD format",
		},
		{
			name: "Transaction with rounding up amount",
			transaction: Transaction{
				Description:     "Round up",
				Amount:          1.006,
				TransactionDate: "2020-01-01",
			},
			expectedResult: "",
		},
		{
			name: "Valid transaction",
			transaction: Transaction{
				Description:     "This is exactly 50 characters long. Juuuust right.",
				Amount:          0.01,
				TransactionDate: "2020-01-01",
			},
			expectedResult: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.transaction.Validate()
			if result != tt.expectedResult {
				t.Errorf("Validate() = %v, want %v", result, tt.expectedResult)
			}
		})
	}
}

func TestRoundToCents(t *testing.T) {
	tests := []struct {
		name           string
		amount         float64
		expectedResult float64
	}{
		{
			name:           "Round up",
			amount:         1.006,
			expectedResult: 1.01,
		},
		{
			name:           "Round down",
			amount:         1.005,
			expectedResult: 1.00,
		},
		{
			name:           "No rounding needed",
			amount:         1.00,
			expectedResult: 1.00,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := roundToCents(tt.amount)
			if got != tt.expectedResult {
				t.Errorf("roundToCents() = %v, want %v", got, tt.expectedResult)
			}
		})
	}
}
