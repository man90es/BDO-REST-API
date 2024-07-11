package handlers

import (
	"net/http"
)

func Catchall(w http.ResponseWriter, r *http.Request) {
	giveBadRequestResponse(w)
}
