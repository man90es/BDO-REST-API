package handlers

import (
	"net/http"
)

func catchall(w http.ResponseWriter, r *http.Request) {
	giveBadRequestResponse(w)
}
