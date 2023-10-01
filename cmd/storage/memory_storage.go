package storage

import (
	"fmt"
	"github.com/msmkdenis/yap-shortener/cmd/utils"
)

type URLStorage struct {
	storage map[string]URL
}

func NewMemoryRepository() Repository {
	return URLStorage{storage: make(map[string]URL)}
}

func (repository URLStorage) Add(u string, host string) URL {
	urlKey := utils.GenerateUniqueURLKey()

	_, ok := repository.storage[urlKey]
	if ok {
		return repository.Add(u, host)
	}

	var url = URL{
		ID:        urlKey,
		Original:  u,
		Shortened: "http://" + host + "/" + urlKey,
	}
	repository.storage[urlKey] = url
	fmt.Println(repository.storage)

	return url
}

func (repository URLStorage) GetByID(key string) (url URL, err error) {
	url, ok := repository.storage[key]
	if !ok {
		return url, fmt.Errorf("URL with id = %s not found", key)
	}

	return url, nil
}
