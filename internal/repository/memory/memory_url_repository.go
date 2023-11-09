package memory

import (
	"github.com/labstack/echo/v4"
	"github.com/msmkdenis/yap-shortener/internal/apperrors"
	"github.com/msmkdenis/yap-shortener/internal/utils"
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

func (r *MemoryURLRepository) Insert(c echo.Context, u model.URL) (*model.URL, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var url = u
	r.storage[u.ID] = u

	return &url, nil
}

func (r *MemoryURLRepository) SelectByID(c echo.Context, key string) (*model.URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	url, ok := r.storage[key]
	if !ok {
		return &url, apperrors.ErrorURLNotFound
	}

	return &url, nil
}

func (r *MemoryURLRepository) SelectAll(c echo.Context) ([]model.URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	values := make([]model.URL, 0, len(r.storage))
	for _, v := range r.storage {
		values = append(values, v)
	}

	return values, nil
}

func (r *MemoryURLRepository) DeleteAll(c echo.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	clear(r.storage)
	return nil
}

func (r *MemoryURLRepository) Ping(c echo.Context) error {
	if r.storage == nil {
		return apperrors.NewValueError("storage is not initialized", utils.Caller(), apperrors.ErrorURLNotFound)
	}

	return nil
}

func (r *MemoryURLRepository) InsertBatch(c echo.Context, urls []model.URL) ([]model.URL, error) {
	return nil, nil
}
