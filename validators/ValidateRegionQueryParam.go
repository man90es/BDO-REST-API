package validators

import (
	"fmt"
	"slices"
	"strings"
)

func ValidateRegionQueryParam(query []string) (region string, ok bool, errorMessage string) {
	if 1 > len(query) {
		return "EU", true, ""
	}

	region = strings.ToUpper(query[0])

	// TODO: Add KR region once the translations are ready
	if !slices.Contains([]string{"EU", "NA", "SA"}, region) {
		return region, false, fmt.Sprintf("Region %v is not supported", region)
	}

	return region, true, ""
}
