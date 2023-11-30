package memory

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"

	"github.com/msmkdenis/yap-shortener/internal/apperrors"
	"github.com/msmkdenis/yap-shortener/internal/model"
	"github.com/msmkdenis/yap-shortener/internal/utils"
)

type URLRepository struct {
	mu      sync.RWMutex
	storage map[string]model.URL
	logger  *zap.Logger
}

func NewURLRepository(logger *zap.Logger) *URLRepository {
	return &URLRepository{
		storage: make(map[string]model.URL),
		logger:  logger,
		mu:      sync.RWMutex{},
	}
}

func (r *URLRepository) DeleteURLByUserID(ctx context.Context, userID string, shortURL string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if url, ok := r.storage[shortURL]; ok {
		if url.UserID == userID {
			url.DeletedFlag = true
			r.storage[shortURL] = url
		}
	}

	return nil
}

func (r *URLRepository) SelectAllByUserID(ctx context.Context, userID string) ([]model.URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	urls := make([]model.URL, 0)
	for _, url := range r.storage {
		if url.UserID == userID {
			urls = append(urls, url)
		}
	}

	if len(urls) == 0 {
		return nil, apperrors.NewValueError(fmt.Sprintf("urls not found by user %s", userID), utils.Caller(), apperrors.ErrURLNotFound)
	}

	return urls, nil
}

func (r *URLRepository) Insert(ctx context.Context, u model.URL) (*model.URL, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	url := u
	r.storage[u.ID] = u

	return &url, nil
}

func (r *URLRepository) SelectByID(ctx context.Context, key string) (*model.URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	url, ok := r.storage[key]
	if !ok {
		return &url, apperrors.ErrURLNotFound
	}

	return &url, nil
}

func (r *URLRepository) SelectAll(ctx context.Context) ([]model.URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	values := make([]model.URL, 0, len(r.storage))
	for _, v := range r.storage {
		values = append(values, v)
	}

	return values, nil
}

func (r *URLRepository) DeleteAll(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	clear(r.storage)
	return nil
}

func (r *URLRepository) Ping(ctx context.Context) error {
	if r.storage == nil {
		return apperrors.NewValueError("storage is not initialized", utils.Caller(), apperrors.ErrURLNotFound)
	}

	return nil
}

func (r *URLRepository) InsertAllOrUpdate(ctx context.Context, urls []model.URL) ([]model.URL, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, v := range urls {
		r.storage[v.ID] = v
	}

	return urls, nil
}
