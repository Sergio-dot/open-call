package main

import (
	"log"

	"github.com/Sergio-dot/open-call/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	// load .env files
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}
	err = godotenv.Load(".env.db")
	if err != nil {
		log.Fatal("Error loading .env.db file: ", err)
	}

	// start server
	if err := server.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
