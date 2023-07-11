package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"bdo-rest-api/httpServer"
	"bdo-rest-api/scrapers"
)

func main() {
	// Parse flags
	flagCacheCap := flag.Int("cachecap", 1e4, "Cache capacity")
	flagCacheTTL := flag.Int("cachettl", 180, "Cache TTL in minutes")
	flagPort := flag.Int("port", 8001, "Port to catch requests on")
	flagProxy := flag.String("proxy", "", "Open proxy address to make requests to BDO servers")
	flagVerbose := flag.Bool("verbose", false, "Print out additional logs into stdout")
	flag.Parse()

	// Read port from flags and env
	var port string
	if *flagPort == 8001 && len(os.Getenv("PORT")) > 0 {
		port = os.Getenv("PORT")
	} else {
		port = fmt.Sprintf("%v", *flagPort)
	}

	// Read proxies from flags
	var proxies []string
	if len(*flagProxy) > 0 {
		proxies = strings.Fields(*flagProxy)
	} else {
		proxies = strings.Fields(os.Getenv("PROXY"))
	}
	scrapers.PushProxies(proxies...)

	// Set scraper verbosity level according to flag
	scrapers.SetVerbose(*flagVerbose)

	// Print out start info
	configPrintOut := "Configuration:\n" +
		fmt.Sprintf("\tPort:\t\t%v\n", port) +
		fmt.Sprintf("\tProxies:\t%v\n", proxies) +
		fmt.Sprintf("\tVerbosity:\t%v\n", *flagVerbose)

	if httpServer.CacheSupport {
		configPrintOut += fmt.Sprintf("\tCache TTL:\t%v minutes\n", *flagCacheTTL) +
			fmt.Sprintf("\tCache capacity:\t%v\n\n", *flagCacheCap)
	} else {
		configPrintOut += "\tCache:\tUnsupported in this build\n\n"
	}

	fmt.Printf(configPrintOut)

	// Build server
	srv := httpServer.BuildServer(&port, flagCacheTTL, flagCacheCap)

	log.Println("Listening for requests")
	log.Fatal(srv.ListenAndServe())
}
