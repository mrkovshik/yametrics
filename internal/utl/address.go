package utl

import "regexp"

func ValidateAddress(addr string) bool {
	pattern := `^([^:]+):(\d+)$`
	regexpPattern := regexp.MustCompile(pattern)
	return regexpPattern.MatchString(addr)
}
