package handler

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mvfavila/transactions/model"
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
