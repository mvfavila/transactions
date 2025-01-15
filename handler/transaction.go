package handler

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mvfavila/transactions/model"
	"github.com/mvfavila/transactions/service"
	"github.com/mvfavila/transactions/util"
)

// StoreTransactionHandler handles POST /transactions.
// It stores a new transaction in the database. The expected body is a JSON object with fields:
// - description: string
// - amount: float64
// - transaction_date: string in YYYY-MM-DD format
//
// If the request body is invalid, it will return 400 with the error message.
// If the transaction is invalid (i.e. description is too long, amount is not positive, or date is invalid), it will return 400 with the error message.
// If the transaction is successfully stored, it will return 201 with the stored transaction in the response body.
func StoreTransactionHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var transaction model.Transaction
		if err := c.ShouldBindJSON(&transaction); err != nil {
			util.InfoLogger.Println(fmt.Sprintf("transaction refused. StatusCode %d:", http.StatusBadRequest), err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if errMsg := transaction.Validate(); errMsg != "" {
			util.InfoLogger.Println(fmt.Sprintf("transaction refused. StatusCode %d:", http.StatusBadRequest), errMsg)
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
			return
		}

		query := "INSERT INTO transactions (description, amount, transaction_date) VALUES (?, ?, ?)"
		res, err := db.Exec(query, transaction.Description, transaction.Amount, transaction.TransactionDate)
		if err != nil {
			util.ErrorLogger.Println(fmt.Sprintf("failed to store transaction. StatusCode %d:", http.StatusInternalServerError), err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store transaction"})
			return
		}

		id, _ := res.LastInsertId()
		transaction.ID = int(id)
		util.InfoLogger.Println("transaction successfully stored:", id)
		c.JSON(http.StatusCreated, transaction)
	}
}

// RetrievePurchaseTransactionHandler handles GET /transactions/:id/exchange-rate/:country.
// It retrieves a transaction, fetches exchange rates, and calculates the converted amount.
func RetrievePurchaseTransactionHandler(db *sql.DB, client *http.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse parameters
		country := c.Param("country")
		id := c.Param("id")

		if country == "" {
			util.WarningLogger.Println("country parameter is required")
			c.JSON(http.StatusBadRequest, gin.H{"error": "country is required"})
			return
		}

		// Retrieve transaction from database
		var transaction model.Transaction
		query := "SELECT id, description, amount, transaction_date FROM transactions WHERE id = ?"
		err := db.QueryRow(query, id).Scan(&transaction.ID, &transaction.Description, &transaction.Amount, &transaction.TransactionDate)
		if err != nil {
			if err == sql.ErrNoRows {
				util.WarningLogger.Printf("transaction with id %s not found", id)
				c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
			} else {
				util.ErrorLogger.Println("failed to retrieve transaction:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve transaction"})
			}
			return
		}

		util.InfoLogger.Println("successfully retrieved transaction:", &transaction)

		// Fetch exchange rates
		rates, err := service.FetchExchangeRates(client, country, &transaction)
		if err != nil {
			util.ErrorLogger.Println("failed to fetch exchange rates:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch exchange rates"})
			return
		}

		// Check if rates were found
		if len(rates) == 0 {
			util.WarningLogger.Printf("no exchange rate found for country %s", country)
			c.JSON(http.StatusNotFound, gin.H{"error": "no exchange rate found"})
			return
		}

		// Use the most recent exchange rate
		latestRate := rates[0]
		convertedAmount := transaction.Amount * latestRate.ExchangeRate

		// Respond with the result
		response := gin.H{
			"id":               transaction.ID,
			"description":      transaction.Description,
			"transaction_date": transaction.TransactionDate,
			"usd_amount":       transaction.Amount,
			"exchange_rate":    latestRate.ExchangeRate,
			"converted_amount": util.RoundToCents(convertedAmount),
		}
		util.InfoLogger.Println("successfully retrieved transaction with exchange rate:", response)
		c.JSON(http.StatusOK, response)
	}
}
