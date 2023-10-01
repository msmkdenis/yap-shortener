package utils

import (
	"github.com/msmkdenis/yap-shortener/cmd/storage"
	"math/rand"
)

func GenerateUniqueURLKey() string {
	runes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")
	urlKey := make([]rune, 8)
	for i := range urlKey {
		urlKey[i] = runes[rand.Intn(len(runes))]
	}

	// Проверка на уникальность через рекурсию, пока не создастся уникальный ключ
	_, ok := storage.Storage[string(urlKey)]
	if ok {
		return GenerateUniqueURLKey()
	}

	return string(urlKey)
}
