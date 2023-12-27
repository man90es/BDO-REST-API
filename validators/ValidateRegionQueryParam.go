package validators

import (
	"slices"
	"strings"
)

func ValidateRegionQueryParam(query []string) (region string, ok bool) {
	if 1 > len(query) {
		return "EU", true
	}

	region = strings.ToUpper(query[0])

	// TODO: Add KR region once the translations are ready
	return region, slices.Contains([]string{"EU", "NA", "SA"}, region)
}
