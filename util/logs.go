package util

import (
	"io"
	"log"
)

var (
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
)

func InitLogger(output io.Writer) {
	// Initialize loggers
	InfoLogger = log.New(output, "INFO: ", log.Ldate|log.Ltime|log.Lmicroseconds)
	WarningLogger = log.New(output, "WARNING: ", log.Ldate|log.Ltime|log.Lmicroseconds)
	ErrorLogger = log.New(output, "ERROR: ", log.Ldate|log.Ltime|log.Lmicroseconds)
}
