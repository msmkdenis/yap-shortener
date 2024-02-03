package hasher

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
)

func GenerateMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	fullHash := hex.EncodeToString(hash[:])
	encoded := base64.StdEncoding.EncodeToString([]byte(fullHash))

	return encoded[:7]
}
