package util

import "testing"

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
			got := RoundToCents(tt.amount)
			if got != tt.expectedResult {
				t.Errorf("roundToCents() = %v, want %v", got, tt.expectedResult)
			}
		})
	}
}
