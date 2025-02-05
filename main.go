package main

import (
	"flag"
	"os"
	"strconv"
	"strings"
	"time"

	"bdo-rest-api/config"
	"bdo-rest-api/handlers"
	"bdo-rest-api/logger"
)

func main() {
	flagCacheTTL := flag.Int("cachettl", 180, "Cache TTL in minutes")
	flagMaintenanceTTL := flag.Int("maintenancettl", 5, "Allows to limit how frequently scraper can check for maintenance end in minutes")
	flagMongo := flag.String("mongo", "", "MongoDB connection string for loggig")
	flagPort := flag.Int("port", 8001, "Port to catch requests on")
	flagProxy := flag.String("proxy", "", "Open proxy address to make requests to BDO servers")
	flagRateLimit := flag.Int64("ratelimit", 512, "Maximum number of requests per minute per IP")
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

	config.SetCacheTTL(time.Duration(*flagCacheTTL) * time.Minute)
	config.SetMaintenanceStatusTTL(time.Duration(*flagMaintenanceTTL) * time.Minute)
	config.SetMongoDB(*flagMongo)
	config.SetRateLimit(*flagRateLimit)
	config.SetVerbosity(*flagVerbose)

	logger.InitLogger()
	handlers.ListenAndServe()
}
