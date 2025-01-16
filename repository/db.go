package repository

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"github.com/mvfavila/transactions/util"
)

// InitializeDB initializes the database.
func InitializeDB(driver string, source string) *sql.DB {
	if _, err := os.Stat(source); os.IsNotExist(err) {
		file, err := os.Create(source)
		if err != nil {
			util.ErrorLogger.Fatalf("Failed to create database file: %v", err)
		}
		file.Close()
	}

	db, err := sql.Open(driver, source)
	if err != nil {
		util.ErrorLogger.Fatalf("Failed to connect to SQLite: %v", err)
	}

	ApplyMigrations(db)
	return db
}

// ApplyMigrations applies the necessary database migrations to the given
// database connection.
func ApplyMigrations(db *sql.DB) {
	migration := `
		CREATE TABLE IF NOT EXISTS transactions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			description TEXT NOT NULL CHECK(length(description) <= 50),
			amount DECIMAL(10, 2) NOT NULL,
			transaction_date TEXT NOT NULL
		);
	`

	if _, err := db.Exec(migration); err != nil {
		util.ErrorLogger.Fatalf("Failed to apply migrations: %v", err)
	}
}
