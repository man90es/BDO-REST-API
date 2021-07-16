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
	flagProxy := flag.String("proxy", "", "Open proxy address to make requests to BDO servers")
	flagPort := flag.Int("port", 8001, "Port to catch requests on")
	flagCacheTTL := flag.Int("cachettl", 180, "Cache TTL in minutes")
	flag.Parse()

	var port string
	if *flagPort == 8001 && len(os.Getenv("PORT")) > 0 {
		port = os.Getenv("PORT")
	} else {
		port = fmt.Sprintf("%v", *flagPort)
	}

	var proxies []string
	if len(*flagProxy) > 0 {
		proxies = strings.Fields(*flagProxy)
	} else {
		proxies = strings.Fields(os.Getenv("PROXY"))
	}
	scrapers.PushProxies(proxies...)

	fmt.Printf("Used configuration:\n\tProxies:\t%v\n\tPort:\t\t%v\n\tCache TTL:\t%v minutes\n\n", proxies, port, *flagCacheTTL)

	srv := httpServer.Server(&port, flagCacheTTL)

	log.Println("Listening for requests")
	log.Fatal(srv.ListenAndServe())
}
