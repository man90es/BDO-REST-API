package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"bdo-rest-api/config"
	"bdo-rest-api/httpServer"
)

func main() {
	flagCacheTTL := flag.Int("cachettl", 180, "Cache TTL in minutes")
	flagPort := flag.Int("port", 8001, "Port to catch requests on")
	flagProxy := flag.String("proxy", "", "Open proxy address to make requests to BDO servers")
	flagVerbose := flag.Bool("verbose", false, "Print out additional logs into stdout")
	flag.Parse()

	// Read port from flags and env
	if *flagPort == 8001 && len(os.Getenv("PORT")) > 0 {
		port, err := strconv.Atoi(os.Getenv("PORT"))

		if nil != err {
			port = 8001
		}

		config.SetPort(port)
	} else {
		config.SetPort(*flagPort)
	}

	// Read proxies from flags
	if len(*flagProxy) > 0 {
		config.SetProxyList(strings.Fields(*flagProxy))
	} else {
		config.SetProxyList(strings.Fields(os.Getenv("PROXY")))
	}

	// Set config variables
	config.SetCacheTTL(time.Duration(*flagCacheTTL) * time.Minute)
	config.SetVerbosity(*flagVerbose)

	config.PrintConfig()
	log.Println("Listening for requests")
	srv := httpServer.BuildServer()
	log.Fatal(srv.ListenAndServe())
}
