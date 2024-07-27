package util

import (
	"errors"
	"strconv"
	"strings"
)

// CutSeconds converts time configurations in format like "2s" to it int representation
// only format with seconds like "xxs" is valid
func CutSeconds(s string) (int, error) {
	reportInterval, found := strings.CutSuffix(s, "s")
	if !found {
		return 0, errors.New("invalid report interval format")
	}
	intReportInterval, err := strconv.Atoi(reportInterval)
	if err != nil {
		return 0, errors.New("invalid report interval format")
	}
	return intReportInterval, nil
}
