package utils

import (
	"crypto/rand"
	"github.com/dchest/authcookie"
	"time"
)

func GenerateToken(id string) string {
	secret := generateSalt(40)
	return authcookie.NewSinceNow(id, 24*time.Hour, secret)
}

func generateSalt(chars int) (salt []byte) {
	saltBytes := make([]byte, chars)
	nRead, err := rand.Read(saltBytes)
	if err != nil {
		salt = []byte{}
	} else if nRead < chars {
		salt = []byte{}
	} else {
		salt = saltBytes
	}
	return
}
