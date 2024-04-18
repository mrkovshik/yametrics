package signature

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"github.com/mrkovshik/yametrics/internal/service"
)

type sha256Sig struct {
	key  string
	body []byte
}

func NewSha256Sig(key string, body []byte) service.Signature {
	return &sha256Sig{
		key:  key,
		body: body,
	}
}

func (s *sha256Sig) Generate() (string, error) {
	h := hmac.New(sha256.New, []byte(s.key))
	if _, err := h.Write(s.body); err != nil {
		return "", err
	}
	dst := h.Sum(nil)
	return hex.EncodeToString(dst), nil
}
