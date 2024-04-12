package retriable

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgerrcode"
)

func Wrap(f func() error) error {
	var (
		retryIntervals = []int{1, 3, 5}
		resultErr      error
	)
	for i := 0; i <= len(retryIntervals); i++ {
		err := f()
		var sysErr *os.SyscallError
		if err == nil || !errors.As(err, &sysErr) {
			return err
		}
		if pgerrcode.IsConnectionException(sysErr.Err.Error()) {
			return err
		}
		if i == len(retryIntervals) {
			return resultErr
		}
		resultErr = errors.Join(resultErr, fmt.Errorf("failed: %v\n retry in %v seconds\n", err, retryIntervals[i]))
		time.Sleep(time.Duration(retryIntervals[i]) * time.Second)
	}
	return resultErr
}
