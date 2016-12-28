package main

import (
	"github.com/goingfullstack/currencyconverter/server"
	"log"
)

func main() {
	// create a new server
	s, err := server.New()
	if err != nil {
		// creation failed, print error and exit
		log.Println("Error creating server:", err)
		return
	}

	err = s.Run() // run the server
	if err != nil {
		// running returned error
		log.Println("Server stopped with error:", err)
	}
}
