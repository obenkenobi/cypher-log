package randutils

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandom32Bytes() (string, error) {
	return GenerateRandomBytes(32)
}

func GenerateRandomBytes(size uint32) (string, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}
