package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/mvfavila/transactions/model"
	"github.com/mvfavila/transactions/repository"
	"github.com/mvfavila/transactions/util"
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

func TestStoreTransactionHandler(t *testing.T) {
	var buf bytes.Buffer

	// Initialize logger with in-memory buffer
	util.InitLogger(&buf)

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to connect to SQLite: %v", err)
	}

	repository.ApplyMigrations(db)

	transaction := model.Transaction{
		Description:     "Test",
		Amount:          1.00,
		TransactionDate: "2020-01-01",
	}

	jsonData, err := json.Marshal(transaction)
	assert.NoError(t, err)

	router := gin.New()
	router.POST("/transactions", StoreTransactionHandler(db))

	req, err := http.NewRequest("POST", "/transactions", strings.NewReader(string(jsonData)))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestRetrievePurchaseTransactionHandler(t *testing.T) {
	t.Run("transaction not found", func(t *testing.T) {
		var buf bytes.Buffer

		// Initialize logger with in-memory buffer
		util.InitLogger(&buf)

		// Initialize the mock database
		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.ExpectQuery("SELECT id, description, amount, transaction_date FROM transactions WHERE id = \\?").
			WithArgs("123").
			WillReturnRows(sqlmock.NewRows([]string{"id", "description", "amount", "transaction_date"}))

		router := gin.New()
		router.GET("/transactions/:id/exchange-rate/:country", RetrievePurchaseTransactionHandler(db, nil))
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/transactions/123/exchange-rate/USD", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, "{\"error\":\"transaction not found\"}", w.Body.String())
	})

	t.Run("exchange rate not found", func(t *testing.T) {
		var buf bytes.Buffer

		// Initialize logger with in-memory buffer
		util.InitLogger(&buf)

		// Initialize the mock database
		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mockResponseBody := `{
			"data": [],
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

		mock.ExpectQuery("SELECT id, description, amount, transaction_date FROM transactions WHERE id = \\?").
			WithArgs("123").
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "description", "amount", "transaction_date"}).
					AddRow(123, "test", 12.34, "2020-01-01"),
			)

		router := gin.New()
		router.GET("/transactions/:id/exchange-rate/:country", RetrievePurchaseTransactionHandler(db, mockClient))
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/transactions/123/exchange-rate/USD", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, "{\"error\":\"no exchange rate found\"}", w.Body.String())
	})

	t.Run("success", func(t *testing.T) {
		var buf bytes.Buffer

		// Initialize logger with in-memory buffer
		util.InitLogger(&buf)

		// Initialize the mock database
		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.ExpectQuery("SELECT id, description, amount, transaction_date FROM transactions WHERE id = \\?").
			WithArgs("123").
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "description", "amount", "transaction_date"}).
					AddRow(123, "test", 12.34, "2020-01-01"),
			)

		mockResponseBody := `{
			"data": [
				{
					"record_date": "2020-01-01",
					"country": "Afghanistan",
					"currency": "Afghani",
					"country_currency_desc": "Afghanistan-Afghani",
					"exchange_rate": "1.0",
					"effective_date": "2020-01-01",
					"src_line_nbr": "1",
					"record_fiscal_year": "2020",
					"record_fiscal_quarter": "1",
					"record_calendar_year": "2020",
					"record_calendar_quarter": "1",
					"record_calendar_month": "1",
					"record_calendar_day": "1"
				}
			],
			"meta": {},
			"links": {}
		}`

		client := &http.Client{
			Transport: &mockRoundTripper{
				mockResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(mockResponseBody)),
				},
			},
		}

		router := gin.New()
		router.GET("/transactions/:id/exchange-rate/:country", RetrievePurchaseTransactionHandler(db, client))
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/transactions/123/exchange-rate/USD", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "{\"converted_amount\":12.34,\"description\":\"test\",\"exchange_rate\":1,\"id\":123,\"transaction_date\":\"2020-01-01\",\"usd_amount\":12.34}", w.Body.String())
	})
}
