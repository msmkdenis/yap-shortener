package memory

import (
	"fmt"
	"github.com/msmkdenis/yap-shortener/internal/domain"
	"github.com/msmkdenis/yap-shortener/internal/utils"
)

type MemoryURLRepository struct {
	storage map[string]domain.URL
}

func NewURLRepository() *MemoryURLRepository {
	return &MemoryURLRepository{storage: make(map[string]domain.URL)}
}

func (r *MemoryURLRepository) Insert(u string, host string) domain.URL {
	urlKey := utils.GenerateMD5Hash(u)

	var url = domain.URL{
		ID:        urlKey,
		Original:  u,
		Shortened: host + "/" + urlKey,
	}
	r.storage[urlKey] = url

	return url
}

func (r *MemoryURLRepository) SelectAll() []string {
	values := make([]string, 0, len(r.storage))
	for _, v := range r.storage {
		values = append(values, v.Original)
	}
	return values
}

func (r *MemoryURLRepository) DeleteAll() {
	clear(r.storage)
}

func (r *MemoryURLRepository) SelectByID(key string) (url domain.URL, err error) {
	url, ok := r.storage[key]
	if !ok {
		return url, fmt.Errorf("URL with id = %s not found", key)
	}
	return url, nil
}
