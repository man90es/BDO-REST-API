//go:build cacheless

package httpServer

import (
	"net/http"
	"time"

	"bdo-rest-api/middleware"

	"github.com/gorilla/mux"
)

const CacheSupport = false

func registerHandlers(handlerMap map[string]func(http.ResponseWriter, *http.Request), ttl time.Duration, cap int) (*mux.Router, error) {
	router := mux.NewRouter()

	for route, handler := range handlerMap {
		router.Handle(route,
			middleware.SetHeaders(
				middleware.CheckForMaintenance(
					http.HandlerFunc(handler),
				),
			),
		).Methods("GET")
	}

	return router, nil
}
