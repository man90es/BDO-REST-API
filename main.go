package main

import (
	"log"
	"os"
	"strings"

	"gitlab.com/man90/black-desert-social-rest-api/httpServer"
	"gitlab.com/man90/black-desert-social-rest-api/scraper"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) < 1 {
		port = "8001"
	}

	proxies := strings.Fields(os.Getenv("PROXY"))
	scraper.PushProxies(proxies...)

	srv := httpServer.Server(&port)

	log.Printf("Listening for requests on port %v.", port)
	log.Fatal(srv.ListenAndServe())
}
