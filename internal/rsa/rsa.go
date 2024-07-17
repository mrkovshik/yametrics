package rsa

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
	"os"

	"github.com/mrkovshik/yametrics/internal/util/retriable"
)

// Encrypt encrypts data using the given RSA public key in PEM format
func Encrypt(publicKeyPem []byte, data []byte) (string, error) {
	// Decode the PEM formatted public key
	block, _ := pem.Decode(publicKeyPem)
	if block == nil {
		return "", errors.New("failed to decode PEM block containing the public key")
	}

	// Parse the public key
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// Type assert the public key to an rsa.PublicKey
	publicKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return "", errors.New("not an RSA public key")
	}

	// Encrypt the data with the public key
	encryptedBytes, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, data)
	if err != nil {
		return "", err
	}

	// Encode the encrypted data in base64 for safe transmission
	encryptedBase64 := base64.StdEncoding.EncodeToString(encryptedBytes)

	return encryptedBase64, nil
}

func Decrypt(privateKeyPem []byte, data []byte) ([]byte, error) {

	block, _ := pem.Decode(privateKeyPem)
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing the private key")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	castedPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("not an RSA private key")
	}

	ciphertext := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	n, err := base64.StdEncoding.Decode(ciphertext, data)
	if err != nil {
		return nil, err
	}
	ciphertext = ciphertext[:n]

	secret, err := rsa.DecryptPKCS1v15(nil, castedPrivateKey, ciphertext)
	if err != nil {
		return nil, err
	}

	return secret, nil

}

// ReadPEMFile reads the PEM file from the given path and returns its contents as a string
func ReadPEMFile(path string) ([]byte, error) {
	file, err := retriable.OpenRetryable(func() (*os.File, error) {
		return os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
	})
	if err != nil {
		return nil, err
	}
	defer file.Close() //nolint:all
	reader := bufio.NewReader(file)
	pemBytes, err := io.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	return pemBytes, nil
}
