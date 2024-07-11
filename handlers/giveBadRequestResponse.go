package handlers

import (
	"encoding/json"
	"net/http"
)

const docsLink = "https://man90es.github.io/BDO-REST-API"

func giveBadRequestResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Route or parameter is invalid. See documentation " + docsLink,
	})
}
