package validators

func ValidateBypassCacheQueryParam(query []string) (bypassCache bool) {
	return len(query) > 0 && query[0] == "1"
}
