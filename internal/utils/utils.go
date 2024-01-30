package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"path"
	"path/filepath"
	"runtime"
)

// GenerateMD5Hash generates MD5 hash from string to shorten URL
func GenerateMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	fullHash := hex.EncodeToString(hash[:])
	encoded := base64.StdEncoding.EncodeToString([]byte(fullHash))

	return encoded[:7]
}

// Caller returns file name and line number of function call
func Caller() string {
	_, file, lineNo, ok := runtime.Caller(1)
	if !ok {
		return "runtime.Caller() failed"
	}

	fileName := path.Base(file)
	dir := filepath.Base(filepath.Dir(file))
	return fmt.Sprintf("%s/%s:%d", dir, fileName, lineNo)
}
