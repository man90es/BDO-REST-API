// +build cacheless

package httpServer

import (
	"fmt"
	"net/http"
	"time"

	"bdo-rest-api/handlers"

	"github.com/gorilla/mux"
)

func Server(port *string, flagCacheTTL *int) (srv *http.Server) {
	router := mux.NewRouter()

	router.HandleFunc("/v1/adventurer/search", handlers.GetAdventurerSearch).Methods("GET")
	router.HandleFunc("/v1/guild/search", handlers.GetGuildSearch).Methods("GET")
	router.HandleFunc("/v1/adventurer", handlers.GetAdventurer).Methods("GET")
	router.HandleFunc("/v1/guild", handlers.GetGuild).Methods("GET")

	srv = &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("0.0.0.0:%v", *port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return
}
