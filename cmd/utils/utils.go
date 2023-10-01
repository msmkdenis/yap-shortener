package utils

import (
	"math/rand"
)

func GenerateUniqueURLKey() string {
	runes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")
	urlKey := make([]rune, 8)
	for i := range urlKey {
		urlKey[i] = runes[rand.Intn(len(runes))]
	}

	return string(urlKey)
}
