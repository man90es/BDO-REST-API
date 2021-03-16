package main

import (
	"log"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) < 1 {
		port = "8001"
	}

	srv := serverFactory(&port)

	log.Printf("Listening for requests on port %v.", port)
	log.Fatal(srv.ListenAndServe())
}
