package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/mvfavila/transactions/config"
	"github.com/mvfavila/transactions/model"
)

// TreasuryRate represents a single exchange rate entry
type TreasuryRate struct {
	Currency      string  `json:"currency"`
	Country       string  `json:"country"`
	ExchangeRate  float64 `json:"exchange_rate,string"`
	EffectiveDate string  `json:"effective_date"`
}

// TreasuryResponse represents the API response structure
type TreasuryResponse struct {
	Data []TreasuryRate `json:"data"`
}

// FetchExchangeRates fetches exchange rates from the Treasury API
func FetchExchangeRates(client *http.Client, country string, transaction *model.Transaction) ([]TreasuryRate, error) {
	if country == "" {
		return nil, fmt.Errorf("country is required")
	}

	var query, err = getRequestQuery(country, transaction)
	if err != nil {
		return nil, err
	}

	// Make an HTTP GET request
	resp, err := client.Get(fmt.Sprintf("%s?filter=%s", config.AppConfig.TreasuryAPIBaseURL, query))
	if err != nil {
		return nil, fmt.Errorf("failed to make request to Treasury API: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response status code is successful
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse JSON response
	var treasuryResponse TreasuryResponse
	if err := json.Unmarshal(body, &treasuryResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}

	return treasuryResponse.Data, nil
}

// getDateMinusSixMonths calculates the date that is six months prior to the given currentDate.
func getDateMinusSixMonths(currentDate string) (string, error) {
	if d, err := time.Parse(config.AppConfig.ExpectedDateFormat, currentDate); err == nil {
		return d.AddDate(0, -6, 0).Format(config.AppConfig.ExpectedDateFormat), nil
	}

	return "", fmt.Errorf("transaction date must be in YYYY-MM-DD format")
}

// getRequestQuery generates a filter and sort string for the Treasury API based on the given country and transaction.
//
// The filter string is based on the following criteria:
// - country: the country of the transaction
// - effective_date: the date of the transaction or the date 6 months prior to the transaction date, whichever is later.
//
// If the transaction date is invalid, an error is returned.
func getRequestQuery(country string, transaction *model.Transaction) (string, error) {
	var bottomDate string
	var err error
	if bottomDate, err = getDateMinusSixMonths(transaction.TransactionDate); err != nil {
		return "", err
	}

	return fmt.Sprintf("country:eq:%s,effective_date:gte:%s,effective_date:lte:%s&sort=-effective_date", country, bottomDate, transaction.TransactionDate), nil
}
