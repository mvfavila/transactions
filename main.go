package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/mvfavila/transactions/handler"
	"github.com/mvfavila/transactions/repository"
	"github.com/mvfavila/transactions/util"
)

const (
	logFile          = "transactions.log"
	port             = ":8080"
	transactionsPath = "/transactions"
)

func main() {
	// Open or create the log file
	logFile, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	// Initialize the logger
	util.InitLogger(logFile)

	// Initialize the database
	db := repository.InitializeDB()
	defer db.Close()

	// Initialize the router
	router := gin.Default()

	router.POST(transactionsPath, handler.StoreTransactionHandler(db))

	util.InfoLogger.Println("transactions service listening on port", port)
	router.Run(port)
}
