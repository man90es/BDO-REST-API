// +build cacheless

package httpServer

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"bdo-rest-api/api"
)

func Server(port *string, flagCacheTTL *int) (srv *http.Server) {
	router := mux.NewRouter()

	router.HandleFunc("/v1/adventurer/search", api.ProfileSearch).Methods("GET")
	router.HandleFunc("/v1/guild/search", api.GuildProfileSearch).Methods("GET")
	router.HandleFunc("/v1/adventurer", api.Profile).Methods("GET")
	router.HandleFunc("/v1/guild", api.GuildProfile).Methods("GET")

	srv = &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("0.0.0.0:%v", *port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return
}
