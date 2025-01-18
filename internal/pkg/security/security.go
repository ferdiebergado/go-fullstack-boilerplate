package security

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomBytes(length uint32) ([]byte, error) {
	key := make([]byte, length)

	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func GenerateRandomBytesEncoded(length uint32) (string, error) {
	key, err := GenerateRandomBytes(length)

	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(key), nil
}
