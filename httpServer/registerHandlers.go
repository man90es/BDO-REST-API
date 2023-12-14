package httpServer

import (
	"net/http"

	"bdo-rest-api/middleware"

	"github.com/gorilla/mux"
)

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
