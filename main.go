package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/mvfavila/transactions/config"
	"github.com/mvfavila/transactions/handler"
	"github.com/mvfavila/transactions/middleware"
	"github.com/mvfavila/transactions/repository"
	"github.com/mvfavila/transactions/util"
)

const transactionsPath = "/transactions"

func main() {
	// Determine the environment
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "prod" // Default to prod if APP_ENV is not set
	}

	if env != "dev" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Load configuration
	if err := config.LoadConfig(env); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	appConfig := config.AppConfig

	// Open or create the log file
	logFile, err := os.OpenFile(appConfig.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
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
	router := gin.New()

	// Attach middleware
	router = middleware.Attach(router)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	router.POST(transactionsPath, handler.StoreTransactionHandler(db))
	router.GET(transactionsPath+"/:id/exchange-rate/:country", handler.RetrievePurchaseTransactionHandler(db, httpClient))

	// Start the application
	util.InfoLogger.Println("transactions service listening on port", appConfig.Port)
	if err := router.Run(":" + appConfig.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
