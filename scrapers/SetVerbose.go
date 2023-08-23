package scrapers

var verbose = false

func SetVerbose(v bool) {
	verbose = v
}

func isVerbose() bool {
	return verbose
}
