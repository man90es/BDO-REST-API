package main

import (
	"flag"
	"os"
	"strconv"
	"strings"
	"time"

	"bdo-rest-api/handlers"
	"bdo-rest-api/logger"
	"bdo-rest-api/scraper"

	"github.com/spf13/viper"
)

func main() {
	flagCacheTTL := flag.Uint("cachettl", 180, "Cache TTL in minutes")
	flagMaintenanceTTL := flag.Uint("maintenancettl", 5, "Allows to limit how frequently scraper can check for maintenance end in minutes")
	flagMongo := flag.String("mongo", "", "MongoDB connection string for loggig")
	flagPort := flag.Uint("port", 8001, "Port to catch requests on")
	flagProxy := flag.String("proxy", "", "Open proxy address to make requests to BDO servers")
	flagRateLimit := flag.Uint64("ratelimit", 512, "Maximum number of requests per minute per IP")
	flagVerbose := flag.Bool("verbose", false, "Print out additional logs into stdout")
	flagTaskRetries := flag.Uint("taskretries", 3, "Number of retries for a scraping task")
	flagMaxTasksPerClient := flag.Uint("maxtasksperclient", 5, "Maximum number of scraping tasks per client")
	flag.Parse()

	// Read port from flags and env
	if *flagPort == 8001 && len(os.Getenv("PORT")) > 0 {
		port, err := strconv.Atoi(os.Getenv("PORT"))

		if err != nil {
			port = 8001
		}

		viper.Set("port", port)
	} else {
		viper.Set("port", int(*flagPort))
	}

	// Read proxies from flags
	if len(*flagProxy) > 0 {
		viper.Set("proxy", strings.Fields(*flagProxy))
	} else {
		viper.Set("proxy", strings.Fields(os.Getenv("PROXY")))
	}

	viper.Set("cachettl", time.Duration(*flagCacheTTL)*time.Minute)
	viper.Set("maintenancettl", time.Duration(*flagMaintenanceTTL)*time.Minute)
	viper.Set("maxtasksperclient", int(*flagMaxTasksPerClient))
	viper.Set("mongo", *flagMongo)
	viper.Set("ratelimit", int64(*flagRateLimit))
	viper.Set("taskretries", int(*flagTaskRetries))
	viper.Set("verbose", *flagVerbose)

	logger.InitLogger()
	scraper.InitScraper()
	handlers.ListenAndServe()
}
