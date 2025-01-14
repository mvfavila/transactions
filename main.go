package main

import (
	"github.com/mvfavila/transactions/repository"
	"github.com/mvfavila/transactions/util"
)

func main() {
	// Initialize the logger
	util.InitLogger()

	// Initialize the database
	db := repository.InitializeDB()
	defer db.Close()

	util.InfoLogger.Println("transactions service started")
	println("Hello World")
}
