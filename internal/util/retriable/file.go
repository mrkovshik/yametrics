package retriable

import (
	"errors"
	"fmt"
	"os"
	"time"
)

// OpenRetryable attempts to open a file with retry logic for handling file system errors.
// Parameters:
// - f: the function to be executed with retry logic that returns a file and an error.
// Returns:
// - a file pointer if successful, otherwise returns an error if the function fails after the specified retries.
func OpenRetryable(f func() (*os.File, error)) (*os.File, error) {
	var (
		retryIntervals = []int{1, 3, 5} // Retry intervals in seconds
		resultErr, err error
		file           *os.File
	)
	for i := 0; i <= len(retryIntervals); i++ {
		file, err = f()
		var fileErr *os.PathError
		if err == nil || !errors.As(err, &fileErr) {
			return file, err
		}
		if i == len(retryIntervals) {
			return file, resultErr
		}
		resultErr = errors.Join(resultErr, fmt.Errorf("failed: %v\n retry in %v seconds", err, retryIntervals[i]))
		time.Sleep(time.Duration(retryIntervals[i]) * time.Second)
	}
	return file, resultErr
}
