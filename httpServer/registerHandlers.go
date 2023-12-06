//go:build !cacheless

package httpServer

import (
	"net/http"
	"time"

	"bdo-rest-api/middleware"

	"github.com/gorilla/mux"
	cache "github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/memory"
)

const CacheSupport = true

func registerHandlers(handlerMap map[string]func(http.ResponseWriter, *http.Request), ttl time.Duration, cap int) (*mux.Router, error) {
	memcached, err := memory.NewAdapter(
		memory.AdapterWithAlgorithm(memory.LRU),
		memory.AdapterWithCapacity(cap),
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
					middleware.CheckForMaintenance(
						http.HandlerFunc(handler),
					),
				),
			),
		).Methods("GET")
	}

	return router, nil
}
