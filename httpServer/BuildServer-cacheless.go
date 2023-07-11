//go:build cacheless

package httpServer

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"bdo-rest-api/handlers"
	"bdo-rest-api/middleware"

	"github.com/gorilla/mux"
)

const CacheSupport = false

func registerHandlers(handlerMap map[string]func(http.ResponseWriter, *http.Request)) (*mux.Router, error) {
	router := mux.NewRouter()

	for route, handler := range handlerMap {
		router.Handle(route,
			middleware.SetHeaders(
				http.HandlerFunc(handler),
			),
		).Methods("GET")
	}

	return router, nil
}

func BuildServer(port *string, flagCacheTTL *int, flagCacheCap *int) (srv *http.Server) {
	router, err := registerHandlers(map[string]func(http.ResponseWriter, *http.Request){
		"/v1/adventurer/search": handlers.GetAdventurerSearch,
		"/v1/guild/search":      handlers.GetGuildSearch,
		"/v1/adventurer":        handlers.GetAdventurer,
		"/v1/guild":             handlers.GetGuild,
	})

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	srv = &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%v", *port),
		Handler:      router,
		IdleTimeout:  60 * time.Second,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	return
}
