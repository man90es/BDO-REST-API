package httpServer

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"bdo-rest-api/config"
	"bdo-rest-api/handlers"
)

func BuildServer() *http.Server {
	router, err := registerHandlers(map[string]func(http.ResponseWriter, *http.Request){
		"/v1":                   handlers.GetStatus,
		"/v1/adventurer":        handlers.GetAdventurer,
		"/v1/adventurer/search": handlers.GetAdventurerSearch,
		"/v1/guild":             handlers.GetGuild,
		"/v1/guild/search":      handlers.GetGuildSearch,
	}, handlers.Catchall)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%v", config.GetPort()),
		Handler:      router,
		IdleTimeout:  60 * time.Second,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
}
