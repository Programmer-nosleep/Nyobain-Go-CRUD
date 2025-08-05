package config

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomKey() (string, error) {
	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)

	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}