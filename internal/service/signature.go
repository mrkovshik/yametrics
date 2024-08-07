// Package service provides interfaces and utilities for generating signatures.
package service

// Signature represents an interface for generating signatures.
type Signature interface {
	// Generate generates a signature and returns it as a string.
	Generate() (string, error)
}
