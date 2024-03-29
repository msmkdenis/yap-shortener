// Package service provides URL service.
package service

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/msmkdenis/yap-shortener/internal/dto"
	"github.com/msmkdenis/yap-shortener/internal/model"
	urlErr "github.com/msmkdenis/yap-shortener/internal/urlerr"
	"github.com/msmkdenis/yap-shortener/pkg/apperr"
	"github.com/msmkdenis/yap-shortener/pkg/hashgen"
)

// URLRepository represents URL repository interface.
type URLRepository interface {
	Insert(ctx context.Context, u model.URL) (*model.URL, error)
	InsertAllOrUpdate(ctx context.Context, urls []model.URL) ([]model.URL, error)
	SelectByID(ctx context.Context, key string) (*model.URL, error)
	SelectAll(ctx context.Context) ([]model.URL, error)
	SelectAllByUserID(ctx context.Context, userID string) ([]model.URL, error)
	DeleteAll(ctx context.Context) error
	DeleteURLByUserID(ctx context.Context, userID string, shortURLs string) error
	SelectStats(ctx context.Context) (*model.URLStats, error)
	Ping(ctx context.Context) error
}

// URLUseCase represents implementation of URL service.
type URLUseCase struct {
	repository URLRepository
	logger     *zap.Logger
}

// NewURLService initializes a new URLUseCase with the given URLRepository and logger.
func NewURLService(repository URLRepository, logger *zap.Logger) *URLUseCase {
	return &URLUseCase{
		repository: repository,
		logger:     logger,
	}
}

// GetStats returns URL stats.
func (u *URLUseCase) GetStats(ctx context.Context) (*dto.URLStats, error) {
	stats, err := u.repository.SelectStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s %w", apperr.Caller(), err)
	}

	response := &dto.URLStats{
		Urls:  stats.Urls,
		Users: stats.Users,
	}
	return response, nil
}

// GetAllByUserID returns all URLs by user ID.
func (u *URLUseCase) GetAllByUserID(ctx context.Context, userID string) ([]dto.URLBatchResponseByUserID, error) {
	urls, err := u.repository.SelectAllByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s %w", apperr.Caller(), err)
	}

	response := make([]dto.URLBatchResponseByUserID, len(urls))
	for i, url := range urls {
		response[i] = dto.URLBatchResponseByUserID{
			OriginalURL: url.Original,
			ShortURL:    url.Shortened,
		}
	}

	return response, nil
}

// DeleteURLByUserID deletes URL by user ID.
func (u *URLUseCase) DeleteURLByUserID(ctx context.Context, userID string, shortURL string) error {
	err := u.repository.DeleteURLByUserID(ctx, userID, shortURL)
	if err != nil {
		return fmt.Errorf("%s %w", apperr.Caller(), err)
	}

	return nil
}

// Add adds a new URL.
func (u *URLUseCase) Add(ctx context.Context, s, host string, userID string) (*model.URL, error) {
	urlKey := hashgen.GenerateMD5Hash(s)
	url := &model.URL{
		ID:          urlKey,
		Original:    s,
		Shortened:   host + "/" + urlKey,
		UserID:      userID,
		DeletedFlag: false,
	}

	existingURL, err := u.repository.SelectByID(ctx, urlKey)
	if err == nil {
		return existingURL, fmt.Errorf("%s %w", apperr.Caller(), urlErr.ErrURLAlreadyExists)
	}

	savedURL, err := u.repository.Insert(ctx, *url)
	if err != nil {
		return nil, fmt.Errorf("%s %w", apperr.Caller(), err)
	}

	return savedURL, nil
}

// GetAll returns all URLs.
func (u *URLUseCase) GetAll(ctx context.Context) ([]string, error) {
	urls, err := u.repository.SelectAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s %w", apperr.Caller(), err)
	}

	originalURLs := make([]string, 0, len(urls))
	for _, url := range urls {
		originalURLs = append(originalURLs, url.Original)
	}

	return originalURLs, nil
}

// DeleteAll deletes all URLs.
func (u *URLUseCase) DeleteAll(ctx context.Context) error {
	if err := u.repository.DeleteAll(ctx); err != nil {
		return fmt.Errorf("%s %w", apperr.Caller(), err)
	}
	return nil
}

// GetByyID returns URL by ID.
func (u *URLUseCase) GetByyID(ctx context.Context, key string) (string, error) {
	url, err := u.repository.SelectByID(ctx, key)
	if err != nil {
		return "", fmt.Errorf("%s %w", apperr.Caller(), err)
	}

	if url.DeletedFlag {
		return "", apperr.NewValueError("deleted url", apperr.Caller(), urlErr.ErrURLDeleted)
	}

	return url.Original, nil
}

// Ping pings the URL repository.
func (u *URLUseCase) Ping(ctx context.Context) error {
	err := u.repository.Ping(ctx)
	return err
}

// AddAll adds URLs.
func (u *URLUseCase) AddAll(ctx context.Context, urls []dto.URLBatchRequest, host string, userID string) ([]dto.URLBatchResponse, error) {
	urlsToSave := make([]model.URL, 0, len(urls))
	keys := make(map[string]string, len(urls))
	for _, v := range urls {
		if _, ok := keys[v.CorrelationID]; ok {
			return nil, apperr.NewValueError("duplicated keys", apperr.Caller(), urlErr.ErrDuplicatedKeys)
		}
		keys[v.CorrelationID] = v.CorrelationID
		shortURL := hashgen.GenerateMD5Hash(v.OriginalURL)
		url := model.URL{
			ID:            shortURL,
			Original:      v.OriginalURL,
			Shortened:     host + "/" + shortURL,
			CorrelationID: v.CorrelationID,
			UserID:        userID,
			DeletedFlag:   false,
		}
		urlsToSave = append(urlsToSave, url)
	}

	savedURLs, err := u.repository.InsertAllOrUpdate(ctx, urlsToSave)
	if err != nil {
		return nil, fmt.Errorf("%s %w", apperr.Caller(), err)
	}

	response := make([]dto.URLBatchResponse, 0, len(savedURLs))
	for _, url := range savedURLs {
		responseURL := dto.URLBatchResponse{
			CorrelationID: url.CorrelationID,
			ShortenedURL:  url.Shortened,
		}
		response = append(response, responseURL)
	}

	return response, nil
}
