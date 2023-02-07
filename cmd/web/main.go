package main

import (
	"github.com/Sergio-dot/open-call/internal/server"
	"log"
)

func main() {
	// start server
	if err := server.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
