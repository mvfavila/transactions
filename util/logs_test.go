package util

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	var buf bytes.Buffer

	// Initialize logger with in-memory buffer
	InitLogger(&buf)

	// Log some messages
	InfoLogger.Println("Test info message")
	WarningLogger.Println("Test warning message")
	ErrorLogger.Println("Test error message")

	// Check the buffer's content
	logOutput := buf.String()

	assert.Contains(t, logOutput, "Test info message")
	assert.Contains(t, logOutput, "Test warning message")
	assert.Contains(t, logOutput, "Test error message")
}
