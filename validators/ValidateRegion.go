package validators

func ValidateRegion(r *string) bool {
	return *r == "EU" || *r == "NA"
}
