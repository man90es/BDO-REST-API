package handlers

import (
	"fmt"
	"net/http"
	"time"

	"bdo-rest-api/logger"
	"bdo-rest-api/middleware"

	"github.com/spf13/viper"
)

func ListenAndServe() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1", getStatus)
	mux.HandleFunc("GET /v1/adventurer", getAdventurer)
	mux.HandleFunc("GET /v1/adventurer/search", getAdventurerSearch)
	mux.HandleFunc("GET /v1/cache", getCacheSummary)
	mux.HandleFunc("GET /v1/cache/{cacheType}", getCache)
	mux.HandleFunc("GET /v1/guild", getGuild)
	mux.HandleFunc("GET /v1/guild/search", getGuildSearch)
	mux.HandleFunc("/", catchall)

	middlewareStack := middleware.CreateStack(
		middleware.GetSetHeadersMiddleware(),
		middleware.GetRateLimitMiddleware(),
	)

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%v", viper.GetInt("port")),
		Handler:      middlewareStack(mux),
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Error(srv.ListenAndServe().Error())
}
