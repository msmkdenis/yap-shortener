package memory

import (
	"fmt"
	"sync"

	"github.com/msmkdenis/yap-shortener/internal/model"
	"go.uber.org/zap"
)

type MemoryURLRepository struct {
	mu      sync.RWMutex
	storage map[string]model.URL
	logger  *zap.Logger
}

func NewURLRepository(logger *zap.Logger) *MemoryURLRepository {
	return &MemoryURLRepository{
		storage: make(map[string]model.URL),
		logger:  logger,
		mu:      sync.RWMutex{},
	}
}

func (r *MemoryURLRepository) Insert(u model.URL) (*model.URL, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var url = u
	r.storage[u.ID] = u

	return &url, nil
}

func (r *MemoryURLRepository) SelectAll() ([]model.URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	values := make([]model.URL, 0, len(r.storage))
	for _, v := range r.storage {
		values = append(values, v)
	}

	return values, nil
}

func (r *MemoryURLRepository) DeleteAll() {
	r.mu.Lock()
	defer r.mu.Unlock()

	clear(r.storage)
}

func (r *MemoryURLRepository) SelectByID(key string) (*model.URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	url, ok := r.storage[key]
	if !ok {
		return &url, fmt.Errorf("URL with id = %s not found", key)
	}

	return &url, nil
}
