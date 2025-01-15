package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/mvfavila/transactions/model"
	"github.com/mvfavila/transactions/repository"
	"github.com/mvfavila/transactions/util"
)

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
