package util

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomString(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
