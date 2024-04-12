package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashData(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	passHash := hex.EncodeToString(h.Sum(nil))

	return passHash
}
