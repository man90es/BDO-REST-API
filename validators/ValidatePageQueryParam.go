package validators

import "strconv"

func ValidatePageQueryParam(query []string) uint16 {
	if 1 > len(query) {
		return 1
	}

	if page, err := strconv.Atoi(query[0]); nil == err {
		return uint16(max(page, 1))
	}

	return 1
}
