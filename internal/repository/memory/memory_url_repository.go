package memory

import (
	"fmt"
	"github.com/msmkdenis/yap-shortener/internal/model"
)

type MemoryURLRepository struct {
	storage map[string]model.URL
}

func NewURLRepository() *MemoryURLRepository {
	return &MemoryURLRepository{storage: make(map[string]model.URL)}
}

func (r *MemoryURLRepository) Insert(u model.URL) (*model.URL, error) {

	var url = u

	r.storage[u.ID] = u

	return &url, nil
}

func (r *MemoryURLRepository) SelectAll() ([]string, error) {
	values := make([]string, 0, len(r.storage))
	for _, v := range r.storage {
		values = append(values, v.Original)
	}
	return values, nil
}

func (r *MemoryURLRepository) DeleteAll() {
	clear(r.storage)
}

func (r *MemoryURLRepository) SelectByID(key string) (*model.URL, error) {
	url, ok := r.storage[key]
	if !ok {
		return &url, fmt.Errorf("URL with id = %s not found", key)
	}
	return &url, nil
}
