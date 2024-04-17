package hashgen

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
)

// GenerateMD5Hash generates hash from text up to 7 symbols
func GenerateMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	fullHash := hex.EncodeToString(hash[:])
	encoded := base64.StdEncoding.EncodeToString([]byte(fullHash))

	return encoded
}
