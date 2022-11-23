package validators

func ValidateRegion(r *string) bool {
	// TODO: Readd KR region once the translations are ready
	// return *r == "EU" || *r == "NA" || *r == "SA" || *r == "KR"
	return *r == "EU" || *r == "NA" || *r == "SA"
}
