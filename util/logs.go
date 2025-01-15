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

// InitLogger initializes the loggers.
//
// The loggers are expected to be used as follows:
//   - InfoLogger: for logging informational messages
//   - WarningLogger: for logging warning messages
//   - ErrorLogger: for logging error messages
func InitLogger(output io.Writer) {
	InfoLogger = log.New(output, "INFO: ", log.Ldate|log.Ltime|log.Lmicroseconds)
	WarningLogger = log.New(output, "WARNING: ", log.Ldate|log.Ltime|log.Lmicroseconds)
	ErrorLogger = log.New(output, "ERROR: ", log.Ldate|log.Ltime|log.Lmicroseconds)
}
