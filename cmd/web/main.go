package main

import (
	"log"

	"github.com/Sergio-dot/open-call/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	// load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	// start server
	if err := server.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
