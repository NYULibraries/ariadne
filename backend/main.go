package main

import (
	"log"
	"net/http"
)

// Run on port 8080
const appPort = "8080"

func main() {
	router := NewRouter()

	log.Println("Listening on port", appPort)
	log.Fatal(http.ListenAndServe(":"+appPort, router))
}
