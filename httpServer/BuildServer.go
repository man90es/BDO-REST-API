package httpServer

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"bdo-rest-api/handlers"
)

func BuildServer(port *string, flagCacheTTL *int, flagCacheCap *int) *http.Server {
	router, err := registerHandlers(map[string]func(http.ResponseWriter, *http.Request){
		"/v1/adventurer/search": handlers.GetAdventurerSearch,
		"/v1/guild/search":      handlers.GetGuildSearch,
		"/v1/adventurer":        handlers.GetAdventurer,
		"/v1/guild":             handlers.GetGuild,
	}, time.Duration(*flagCacheTTL)*time.Minute, *flagCacheCap)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%v", *port),
		Handler:      router,
		IdleTimeout:  60 * time.Second,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
}
