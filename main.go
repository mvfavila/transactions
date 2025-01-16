package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/mvfavila/transactions/config"
	"github.com/mvfavila/transactions/handler"
	"github.com/mvfavila/transactions/repository"
	"github.com/mvfavila/transactions/util"
)

const transactionsPath = "/transactions"

func main() {
	// Determine the environment
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev" // Default to "dev" environment
	} else if env == "prod" {
		// Set gin to release mode for production
		gin.SetMode(gin.ReleaseMode)
	}

	// Load configuration
	err := config.LoadConfig(env)
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	appConfig := config.AppConfig

	// Open or create the log file
	logFile, err := os.OpenFile(appConfig.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	defer logFile.Close()

	// Initialize the logger
	util.InitLogger(logFile)

	// Initialize the database
	db := repository.InitializeDB(appConfig.Database.Driver, appConfig.Database.Source)
	defer db.Close()

	// Initialize the HTTP client
	httpClient := &http.Client{}

	// Initialize the router
	router := gin.Default()

	router.POST(transactionsPath, handler.StoreTransactionHandler(db))
	router.GET(transactionsPath+"/:id/exchange-rate/:country", handler.RetrievePurchaseTransactionHandler(db, httpClient))

	util.InfoLogger.Println("transactions service listening on port", appConfig.Port)
	router.Run(fmt.Sprintf(":%s", appConfig.Port))
}
