package validators

import "strconv"

func ValidatePage(p *string) bool {
	page, ok := strconv.Atoi(*p)
	return ok == nil && page > 0
}
