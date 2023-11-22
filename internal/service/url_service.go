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

type URLService interface {
	Add(ctx context.Context, s string, host string) (*model.URL, error)
	AddAll(ctx context.Context, urls []dto.URLBatchRequestType, host string) ([]dto.URLBatchResponseType, error)
	GetAll(ctx context.Context) ([]string, error)
	DeleteAll(ctx context.Context) error
	GetByyID(ctx context.Context, key string) (string, error)
	Ping(ctx context.Context) error
}

type URLUseCase struct {
	repository model.URLRepository
	logger     *zap.Logger
}

func NewURLService(repository model.URLRepository, logger *zap.Logger) *URLUseCase {
	return &URLUseCase{
		repository: repository,
		logger:     logger,
	}
}

func (u *URLUseCase) Add(ctx context.Context, s, host string) (*model.URL, error) {
	urlKey := utils.GenerateMD5Hash(s)
	url := &model.URL{
		ID:        urlKey,
		Original:  s,
		Shortened: host + "/" + urlKey,
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

func (u *URLUseCase) AddAll(ctx context.Context, urls []dto.URLBatchRequestType, host string) ([]dto.URLBatchResponseType, error) {
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
		}
		urlsToSave = append(urlsToSave, url)
	}

	savedURLs, err := u.repository.InsertAllOrUpdate(ctx, urlsToSave)
	if err != nil {
		return nil, fmt.Errorf("caller: %s %w", utils.Caller(), err)
	}

	var response []dto.URLBatchResponseType
	for _, url := range savedURLs {
		responseURL := dto.URLBatchResponseType{
			CorrelationID: url.CorrelationID,
			ShortenedURL:  url.Shortened,
		}
		response = append(response, responseURL)
	}

	return response, nil
}
