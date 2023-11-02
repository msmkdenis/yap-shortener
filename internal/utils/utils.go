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

func GenerateMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	fullHash := hex.EncodeToString(hash[:])
	encoded := base64.StdEncoding.EncodeToString([]byte(fullHash))

	return encoded[:7]
}

func Caller() (string) {

	_, file, lineNo, ok := runtime.Caller(1)
	if !ok {
	 return "runtime.Caller() failed"
	}

	fileName := path.Base(file)
	dir := filepath.Base(filepath.Dir(file))
	return fmt.Sprintf("%s/%s:%d", dir, fileName, lineNo)
   }
