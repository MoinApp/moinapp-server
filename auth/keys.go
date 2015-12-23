// Package auth is adapted from https://github.com/alternaDev/georenting-server/blob/ffcba530e65a5469089002bfbaf936d6e9a7fa70/auth/keys.go
package auth

// BUG(sgade): check export necessity

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

var (
	ErrInvalidKeyFormat = errors.New("The key format is invalid.")
)

// GenerateNewPrivateKey generates a new Private Key
func GenerateNewPrivateKey() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)

	return privateKey, err
}

// PrivateKeyToString converts a RSA private Key to a string
func PrivateKeyToString(privateKey *rsa.PrivateKey) string {
	keyBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)

	return string(keyBytes[:])
}

// StringToPrivateKey converts a RSA Private Key String to a private KEy
func StringToPrivateKey(keyString string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(keyString))

	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, ErrInvalidKeyFormat
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)

	if err != nil {
		return nil, err
	}

	return privateKey, nil
}
