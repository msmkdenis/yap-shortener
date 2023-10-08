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
	urlKey := utils.GenerateMD5Hash(u)

	var url = URL{
		ID:        urlKey,
		Original:  u,
		Shortened: "http://" + host + "/" + urlKey,
	}
	repository.storage[urlKey] = url

	return url
}

func (repository URLStorage) GetAll() []string {
	values := make([]string, 0, len(repository.storage))
	for _, v := range repository.storage {
		values = append(values, v.Original)
	}
	return values
}

func (repository URLStorage) DeleteAll() {
	clear(repository.storage)
}

func (repository URLStorage) GetByID(key string) (url URL, err error) {
	url, ok := repository.storage[key]
	if !ok {
		return url, fmt.Errorf("URL with id = %s not found", key)
	}
	return url, nil
}
