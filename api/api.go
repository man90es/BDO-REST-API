package api

type errorResponse struct {
	Error string `json:"error"`
}

func validateRegion(r string) bool {
	return r == "EU" || r == "NA"
}
