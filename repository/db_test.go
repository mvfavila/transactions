package repository

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyMigrations(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to connect to SQLite: %v", err)
	}

	applyMigrations(db)

	// Check if the table exists
	var result string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='transactions'").Scan(&result)
	if err != nil {
		t.Fatalf("Failed to query the database: %v", err)
	}

	assert.Equal(t, "transactions", result)
}
