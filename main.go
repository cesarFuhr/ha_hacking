package main

import (
	"log"
	"net/http"
)

func main() {

	server := NewHTTPServer()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
