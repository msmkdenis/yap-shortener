package service

import (
	"context"
	"fmt"

	"github.com/msmkdenis/yap-shortener/internal/apperrors"
	"github.com/msmkdenis/yap-shortener/internal/handlers/dto"
	"github.com/msmkdenis/yap-shortener/internal/model"
	"github.com/msmkdenis/yap-shortener/internal/utils"
	"go.uber.org/zap"
)

type URLRepository interface {
	Insert(ctx context.Context, u model.URL) (*model.URL, error)
	InsertAllOrUpdate(ctx context.Context, urls []model.URL) ([]model.URL, error)
	SelectByID(ctx context.Context, key string) (*model.URL, error)
	SelectAll(ctx context.Context) ([]model.URL, error)
	SelectAllByUserID(ctx context.Context, userID string) ([]model.URL, error)
	DeleteAll(ctx context.Context) error
	DeleteAllByUserID(ctx context.Context, userID string, shortURLs []string) ([]model.URL, error)
	Ping(ctx context.Context) error
}

type URLUseCase struct {
	repository URLRepository
	logger     *zap.Logger
}

func NewURLService(repository URLRepository, logger *zap.Logger) *URLUseCase {
	return &URLUseCase{
		repository: repository,
		logger:     logger,
	}
}

func (u *URLUseCase) GetAllByUserID(ctx context.Context, userID string) ([]dto.URLBatchResponseByUserID, error) {
	urls, err := u.repository.SelectAllByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("caller: %s %w", utils.Caller(), err)
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

func (u *URLUseCase) DeleteAllByUserID(ctx context.Context, userID string, shortURLs []string) ([]dto.URLBatchResponseByUserID, error) {
	urls, err := u.repository.DeleteAllByUserID(ctx, userID, shortURLs)
	if err != nil {
		return nil, fmt.Errorf("caller: %s %w", utils.Caller(), err)
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

func (u *URLUseCase) Add(ctx context.Context, s, host string, userID string) (*model.URL, error) {
	urlKey := utils.GenerateMD5Hash(s)
	url := &model.URL{
		ID:          urlKey,
		Original:    s,
		Shortened:   host + "/" + urlKey,
		UserID:      userID,
		DeletedFlag: false,
	}

	existingURL, err := u.repository.SelectByID(ctx, urlKey)
	if err == nil {
		return existingURL, apperrors.ErrURLAlreadyExists
	}

	savedURL, err := u.repository.Insert(ctx, *url)
	if err != nil {
		return nil, fmt.Errorf("caller: %s %w", utils.Caller(), err)
	}

	return savedURL, nil
}

func (u *URLUseCase) GetAll(ctx context.Context) ([]string, error) {
	urls, err := u.repository.SelectAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("caller: %s %w", utils.Caller(), err)
	}

	originalURLs := []string{}
	for _, url := range urls {
		originalURLs = append(originalURLs, url.Original)
	}

	return originalURLs, nil
}

func (u *URLUseCase) DeleteAll(ctx context.Context) error {
	if err := u.repository.DeleteAll(ctx); err != nil {
		return fmt.Errorf("caller: %s %w", utils.Caller(), err)
	}
	return nil
}

func (u *URLUseCase) GetByyID(ctx context.Context, key string) (string, error) {
	url, err := u.repository.SelectByID(ctx, key)
	if err != nil {
		return "", fmt.Errorf("caller: %s %w", utils.Caller(), err)
	}

	return url.Original, nil
}

func (u *URLUseCase) Ping(ctx context.Context) error {
	err := u.repository.Ping(ctx)
	return err
}

func (u *URLUseCase) AddAll(ctx context.Context, urls []dto.URLBatchRequest, host string, userID string) ([]dto.URLBatchResponse, error) {
	var urlsToSave []model.URL
	var keys = make(map[string]string, len(urls))
	for _, v := range urls {
		if _, ok := keys[v.CorrelationID]; ok {
			return nil, apperrors.NewValueError("duplicated keys", utils.Caller(), apperrors.ErrDuplicatedKeys)
		}
		keys[v.CorrelationID] = v.CorrelationID
		shortURL := utils.GenerateMD5Hash(v.OriginalURL)
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
		return nil, fmt.Errorf("caller: %s %w", utils.Caller(), err)
	}

	var response []dto.URLBatchResponse
	for _, url := range savedURLs {
		responseURL := dto.URLBatchResponse{
			CorrelationID: url.CorrelationID,
			ShortenedURL:  url.Shortened,
		}
		response = append(response, responseURL)
	}

	return response, nil
}
