package memory

import (
	"fmt"
	"github.com/msmkdenis/yap-shortener/internal/model"
	"go.uber.org/zap"
)

type MemoryURLRepository struct {
	storage map[string]model.URL
	logger  *zap.Logger
}

func NewURLRepository(logger *zap.Logger) *MemoryURLRepository {
	return &MemoryURLRepository{
		storage: make(map[string]model.URL),
		logger:  logger,
	}
}

func (r *MemoryURLRepository) Insert(u model.URL) (*model.URL, error) {
	var url = u
	r.storage[u.ID] = u

	return &url, nil
}

func (r *MemoryURLRepository) SelectAll() ([]model.URL, error) {
	values := make([]model.URL, 0, len(r.storage))
	for _, v := range r.storage {
		values = append(values, v)
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
