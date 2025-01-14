package util

import (
	"log"
	"os"
)

var (
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
)

func InitLogger() {
	// Open or create the log file
	logFile, err := os.OpenFile("transactions.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Initialize loggers
	InfoLogger = log.New(logFile, "INFO: ", log.Ldate|log.Ltime)
	WarningLogger = log.New(logFile, "WARNING: ", log.Ldate|log.Ltime)
	ErrorLogger = log.New(logFile, "ERROR: ", log.Ldate|log.Ltime)
}
