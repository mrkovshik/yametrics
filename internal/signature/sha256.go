// Package signature provides utilities for generating SHA-256 HMAC signatures.
package signature

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"github.com/mrkovshik/yametrics/internal/service"
)

// sha256Sig implements the service.Signature interface for SHA-256 HMAC signatures.
type sha256Sig struct {
	key  string
	body []byte
}

// NewSha256Sig creates a new instance of sha256Sig with the provided key and body.
func NewSha256Sig(key string, body []byte) service.Signature {
	return &sha256Sig{
		key:  key,
		body: body,
	}
}

// Generate generates a SHA-256 HMAC signature for the stored body using the stored key.
func (s *sha256Sig) Generate() (string, error) {
	h := hmac.New(sha256.New, []byte(s.key))
	if _, err := h.Write(s.body); err != nil {
		return "", err
	}
	dst := h.Sum(nil)
	return hex.EncodeToString(dst), nil
}
