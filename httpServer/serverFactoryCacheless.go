// +build cacheless

package httpServer

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"gitlab.com/man90/black-desert-social-rest-api/api"
)

func Server(port *string) (srv *http.Server) {
	router := mux.NewRouter()

	router.HandleFunc("/v0/guildProfile", api.GuildProfile).Methods("GET")
	router.HandleFunc("/v0/profile", api.Profile).Methods("GET")
	router.HandleFunc("/v0/guildProfileSearch", api.GuildProfileSearch).Methods("GET")
	router.HandleFunc("/v0/profileSearch", api.ProfileSearch).Methods("GET")

	srv = &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("0.0.0.0:%v", *port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return
}
