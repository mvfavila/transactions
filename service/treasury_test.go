package service

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/mvfavila/transactions/config"
	"github.com/mvfavila/transactions/model"
	"github.com/stretchr/testify/assert"
)

type mockRoundTripper struct {
	mockResponse *http.Response
	mockError    error
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.mockError != nil {
		return nil, m.mockError
	}
	return m.mockResponse, nil
}

func TestFetchExchangeRates(t *testing.T) {
	// Load default config for testing
	config.LoadDefaultConfig()

	// Mock response body
	mockResponseBody := `{
		"data": [
			{
				"record_date": "2024-12-31",
				"country": "Afghanistan",
				"currency": "Afghani",
				"country_currency_desc": "Afghanistan-Afghani",
				"exchange_rate": "70.35",
				"effective_date": "2024-12-31",
				"src_line_nbr": "1",
				"record_fiscal_year": "2025",
				"record_fiscal_quarter": "1",
				"record_calendar_year": "2024",
				"record_calendar_quarter": "4",
				"record_calendar_month": "12",
				"record_calendar_day": "31"
			}
		],
		"meta": {},
		"links": {}
	}`

	// Create a mock HTTP client
	mockClient := &http.Client{
		Transport: &mockRoundTripper{
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(mockResponseBody)),
			},
		},
	}

	// Call the function
	rates, err := FetchExchangeRates(mockClient, "Brazil", &model.Transaction{TransactionDate: "2025-01-13"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Assert the results
	if len(rates) != 1 {
		t.Fatalf("expected 1 rate, got %d", len(rates))
	}

	expectedRate := 70.35
	if rates[0].ExchangeRate != expectedRate {
		t.Errorf("expected currency name %.2f, got %.2f", expectedRate, rates[0].ExchangeRate)
	}
}

func TestFetchExchangeRates_Error(t *testing.T) {
	// Create a mock HTTP client with an error response
	mockClient := &http.Client{
		Transport: &mockRoundTripper{
			mockError: errors.New("mock error"),
		},
	}

	// Call the function
	_, err := FetchExchangeRates(mockClient, "Brazil", &model.Transaction{TransactionDate: "2025-01-13"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to make request to Treasury API")
	assert.Contains(t, err.Error(), "mock error")
}

func TestGetDateMinusSixMonths(t *testing.T) {
	tests := []struct {
		inputDate       string
		expectedOutput  string
		expectedSuccess bool
	}{
		{
			inputDate:       "2020-01-01",
			expectedOutput:  "2019-07-01",
			expectedSuccess: true,
		},
		{
			inputDate:       "2020-02-29",
			expectedOutput:  "2019-08-29",
			expectedSuccess: true,
		},
		{
			inputDate:       "abc",
			expectedOutput:  "",
			expectedSuccess: false,
		},
		{
			inputDate:       "",
			expectedOutput:  "",
			expectedSuccess: false,
		},
	}

	for _, test := range tests {
		got, err := getDateMinusSixMonths(test.inputDate)
		if test.expectedSuccess {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedOutput, got)
		} else {
			assert.Error(t, err)
			assert.Equal(t, "", got)
		}
	}
}

func TestGetRequestQuery(t *testing.T) {
	var getRequestFilterTests = []struct {
		country         string
		transactionDate string
		expectedFilter  string
		expectedError   error
	}{
		{"Autralia", "2020-01-01", "country:eq:Autralia,effective_date:gte:2019-07-01,effective_date:lte:2020-01-01&sort=-effective_date", nil},
		{"Austria", "2020-02-29", "country:eq:Austria,effective_date:gte:2019-08-29,effective_date:lte:2020-02-29&sort=-effective_date", nil},
		{"Brazil", "abc", "", fmt.Errorf("transaction date must be in YYYY-MM-DD format")},
		{"Czech. Republic", "", "", fmt.Errorf("transaction date must be in YYYY-MM-DD format")},
	}

	for _, tt := range getRequestFilterTests {
		gotFilter, gotErr := getRequestQuery(tt.country, &model.Transaction{TransactionDate: tt.transactionDate})
		if tt.expectedError == nil {
			assert.NoError(t, gotErr)
			assert.Equal(t, tt.expectedFilter, gotFilter)
		} else {
			assert.EqualError(t, gotErr, tt.expectedError.Error())
			assert.Equal(t, "", gotFilter)
		}
	}
}
