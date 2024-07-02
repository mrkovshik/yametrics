// Package retriable provides utilities for executing database operations with retry logic to handle connection exceptions.
package retriable

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgerrcode"
)

// ExecRetryable executes a function with retry logic for handling connection exceptions.
// Parameters:
// - f: the function to be executed with retry logic.
// Returns:
// - an error if the function fails after the specified retries.
func ExecRetryable(f func() error) error {
	var (
		retryIntervals = []int{1, 3, 5} // Retry intervals in seconds
		resultErr      error
	)
	for i := 0; i <= len(retryIntervals); i++ {
		err := f()
		var sysErr *os.SyscallError
		if err == nil || !errors.As(err, &sysErr) {
			return err
		}
		if !pgerrcode.IsConnectionException(sysErr.Err.Error()) {
			return err
		}
		if i == len(retryIntervals) {
			return resultErr
		}
		resultErr = errors.Join(resultErr, fmt.Errorf("failed: %v\n retry in %v seconds", err, retryIntervals[i]))
		time.Sleep(time.Duration(retryIntervals[i]) * time.Second)
	}
	return resultErr
}

func QueryRowRetryable(f func() *sql.Row) (*sql.Row, error) {
	var (
		retryIntervals = []int{1, 3, 5}
		resultErr      error
		row            *sql.Row
	)
	for i := 0; i <= len(retryIntervals); i++ {
		row = f()
		var sysErr *os.SyscallError
		err := row.Err()
		if err == nil || !errors.As(err, &sysErr) {
			return row, err
		}
		if !pgerrcode.IsConnectionException(sysErr.Err.Error()) {
			return row, err
		}
		if i == len(retryIntervals) {
			return row, resultErr
		}
		resultErr = errors.Join(resultErr, fmt.Errorf("failed: %v\n retry in %v seconds", err, retryIntervals[i]))
		time.Sleep(time.Duration(retryIntervals[i]) * time.Second)
	}
	return row, resultErr
}

func QueryRetryable(f func() (*sql.Rows, error)) (*sql.Rows, error) {
	var (
		retryIntervals = []int{1, 3, 5}
		resultErr      error
		rows           *sql.Rows
		err            error
	)
	for i := 0; i <= len(retryIntervals); i++ {
		rows, err = f()
		var sysErr *os.SyscallError
		if err == nil || !errors.As(err, &sysErr) {
			return rows, err
		}
		if !pgerrcode.IsConnectionException(sysErr.Err.Error()) {
			return rows, err
		}
		if i == len(retryIntervals) {
			return rows, resultErr
		}
		resultErr = errors.Join(resultErr, fmt.Errorf("failed: %v\n retry in %v seconds", err, retryIntervals[i]))
		time.Sleep(time.Duration(retryIntervals[i]) * time.Second)
	}
	return rows, resultErr
}
