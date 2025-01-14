package main

import "github.com/mvfavila/transactions/util"

func main() {
	// Initialize the logger
	util.InitLogger()

	util.InfoLogger.Println("transactions service started")
	println("Hello World")
}
