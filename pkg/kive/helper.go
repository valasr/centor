package kive

import (
	"crypto/sha256"
	"encoding/hex"
)

func generateHash(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	hashInBytes := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)
	return hashString
}
