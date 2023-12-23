package helper

import (
	"crypto/sha256"
	"encoding/hex"
)

func Hash(str string) string {
	id := sha256.Sum256([]byte(str))
	return hex.EncodeToString(id[:])[:6]
}
