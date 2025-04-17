package utils

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
)

func GenerateRandomString() (string, error) {
	bytes := make([]byte, 15)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	// Remove 0, O, 1, I to remove ambiguity
	code := base32.NewEncoding("ABCDEFGHJKLMNPQRSTUVWXYZ23456789").EncodeToString(bytes)
	return code, nil
}

func EncodeBase64UrlNoPadding(data []byte) string {
	return base64.RawURLEncoding.WithPadding(base64.NoPadding).EncodeToString(data)
}

func GenerateRandomStringNoPadding() (string, error) {
	bytes := make([]byte, 15)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	code := base32.NewEncoding("ABCDEFGHJKLMNPQRSTUVWXYZ23456789").WithPadding(base32.NoPadding).EncodeToString(bytes)
	return code, nil
}
