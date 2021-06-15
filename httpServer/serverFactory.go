// +build !cacheless

package httpServer

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	cache "github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/memory"

	"black-desert-social-rest-api/api"
)

func Server(port *string, flagCacheTTL *int) (srv *http.Server) {
	memcached, err := memory.NewAdapter(
		memory.AdapterWithAlgorithm(memory.LRU),
		memory.AdapterWithCapacity(1000000),
	)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	cacheClient, err := cache.NewClient(
		cache.ClientWithAdapter(memcached),
		cache.ClientWithTTL(time.Duration(*flagCacheTTL)*time.Minute),
		cache.ClientWithRefreshKey("opn"),
	)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	router := mux.NewRouter()

	router.Handle("/v0/guildProfile", cacheClient.Middleware(http.HandlerFunc(api.GuildProfile))).Methods("GET")
	router.Handle("/v0/profile", cacheClient.Middleware(http.HandlerFunc(api.Profile))).Methods("GET")
	router.Handle("/v0/guildProfileSearch", cacheClient.Middleware(http.HandlerFunc(api.GuildProfileSearch))).Methods("GET")
	router.Handle("/v0/profileSearch", cacheClient.Middleware(http.HandlerFunc(api.ProfileSearch))).Methods("GET")

	srv = &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("0.0.0.0:%v", *port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return
}
