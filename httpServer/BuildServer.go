//go:build !cacheless

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
	cache "github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/memory"
)

func registerHandlers(handlerMap map[string]func(http.ResponseWriter, *http.Request), ttl time.Duration) (*mux.Router, error) {
	memcached, err := memory.NewAdapter(
		memory.AdapterWithAlgorithm(memory.LRU),
		memory.AdapterWithCapacity(1e6),
	)

	if err != nil {
		return nil, err
	}

	cacheClient, err := cache.NewClient(
		cache.ClientWithAdapter(memcached),
		cache.ClientWithTTL(ttl),
		cache.ClientWithRefreshKey("opn"),
		cache.ClientWithStatusCodeFilter(func(code int) bool { return code != 400 }),
		cache.ClientWithExpiresHeader(),
	)

	if err != nil {
		return nil, err
	}

	router := mux.NewRouter()

	for route, handler := range handlerMap {
		router.Handle(route,
			cacheClient.Middleware(
				middleware.SetHeaders(
					http.HandlerFunc(handler),
				),
			),
		).Methods("GET")
	}

	return router, nil
}

func BuildServer(port *string, flagCacheTTL *int) (srv *http.Server) {
	router, err := registerHandlers(map[string]func(http.ResponseWriter, *http.Request){
		"/v1/adventurer/search": handlers.GetAdventurerSearch,
		"/v1/guild/search":      handlers.GetGuildSearch,
		"/v1/adventurer":        handlers.GetAdventurer,
		"/v1/guild":             handlers.GetGuild,
	}, time.Duration(*flagCacheTTL)*time.Minute)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	srv = &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("0.0.0.0:%v", *port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return
}
