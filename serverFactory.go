// +build !cacheless

package main

import (
	"net/http"
	"log"
	"time"
	"fmt"
	"os"

	"github.com/gorilla/mux"
	"github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/memory"

	"gitlab.com/man90/black-desert-social-rest-api/api"
)

func serverFactory(port *string) (srv *http.Server) {
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
		cache.ClientWithTTL(2 * time.Hour),
		cache.ClientWithRefreshKey("opn"),
	)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	router := mux.NewRouter()

	router.Handle("/v0/guildProfile", 		cacheClient.Middleware(http.HandlerFunc(api.GuildProfile)))		.Methods("GET")
	router.Handle("/v0/profile", 			cacheClient.Middleware(http.HandlerFunc(api.Profile)))			.Methods("GET")
	router.Handle("/v0/guildProfileSearch", cacheClient.Middleware(http.HandlerFunc(api.GuildProfileSearch))).Methods("GET")
	router.Handle("/v0/profileSearch", 		cacheClient.Middleware(http.HandlerFunc(api.ProfileSearch)))		.Methods("GET")

	srv = &http.Server{
		Handler: 		router,
		Addr: 			fmt.Sprintf("0.0.0.0:%v", *port),
		WriteTimeout: 	15 * time.Second,
		ReadTimeout:  	15 * time.Second,
	}

	return
}
