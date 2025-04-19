package main

import (
	"flag"
	"os"
	"strconv"
	"strings"
	"time"

	"bdo-rest-api/handlers"
	"bdo-rest-api/logger"

	"github.com/spf13/viper"
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

		if err != nil {
			port = 8001
		}

		viper.Set("port", port)
	} else {
		viper.Set("port", *flagPort)
	}

	// Read proxies from flags
	if len(*flagProxy) > 0 {
		viper.Set("proxy", strings.Fields(*flagProxy))
	} else {
		viper.Set("proxy", strings.Fields(os.Getenv("PROXY")))
	}

	viper.Set("cachettl", time.Duration(*flagCacheTTL)*time.Minute)
	viper.Set("maintenancettl", time.Duration(*flagMaintenanceTTL)*time.Minute)
	viper.Set("mongo", *flagMongo)
	viper.Set("ratelimit", *flagRateLimit)
	viper.Set("verbose", *flagVerbose)

	logger.InitLogger()
	handlers.ListenAndServe()
}
