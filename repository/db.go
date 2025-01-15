package repository

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"github.com/mvfavila/transactions/util"
)

const dbFile = "transactions.db"

// InitializeDB initializes the database.
func InitializeDB() *sql.DB {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		file, err := os.Create(dbFile)
		if err != nil {
			util.ErrorLogger.Fatalf("Failed to create database file: %v", err)
		}
		file.Close()
	}

	db, err := sql.Open("sqlite3", dbFile)
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
