package service

import (
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/msmkdenis/yap-shortener/internal/handlers/dto"
	"github.com/msmkdenis/yap-shortener/internal/model"
	"github.com/msmkdenis/yap-shortener/internal/utils"
	"go.uber.org/zap"
)

type URLService interface {
	Add(ctx echo.Context, s string, host string) (*model.URL, error)
	AddBatch(ctx echo.Context, urls []dto.URLBatchRequestType, host string) ([]model.URL, error)
	GetAll(ctx echo.Context) ([]string, error)
	DeleteAll(ctx echo.Context) error
	GetByyID(ctx echo.Context, key string) (string, error)
	Ping(ctx echo.Context) error
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

func (u *URLUseCase) Add(ctx echo.Context, s, host string) (*model.URL, error) {
	urlKey := utils.GenerateMD5Hash(s)
	url := &model.URL{
		ID:        urlKey,
		Original:  s,
		Shortened: host + "/" + urlKey,
	}

	existingURL, err := u.repository.SelectByID(ctx, url.ID)
	if err == nil {
		return existingURL, nil
	}

	savedURL, err := u.repository.Insert(ctx, *url)
	if err != nil {
		return nil, fmt.Errorf("caller: %s %w", utils.Caller(), err)
	}

	return savedURL, nil
}

func (u *URLUseCase) GetAll(ctx echo.Context) ([]string, error) {
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

func (u *URLUseCase) DeleteAll(ctx echo.Context) error {
	if err := u.repository.DeleteAll(ctx); err != nil {
		return fmt.Errorf("caller: %s %w", utils.Caller(), err)
	}
	return u.repository.DeleteAll(ctx)
}

func (u *URLUseCase) GetByyID(ctx echo.Context, key string) (string, error) {
	url, err := u.repository.SelectByID(ctx, key)
	if err != nil {
		return "", fmt.Errorf("caller: %s %w", utils.Caller(), err)
	}

	return url.Original, nil
}

func (u *URLUseCase) Ping(ctx echo.Context) error {
	err := u.repository.Ping(ctx)
	return err
}

func (u *URLUseCase) AddBatch(ctx echo.Context, urls []dto.URLBatchRequestType, host string) ([]model.URL, error) {
	var urlsToSave []model.URL
	for _, v := range urls {
		shortUrl := utils.GenerateMD5Hash(v.OriginalURL)
		url := model.URL{
			ID:        v.CorrelationID,
			Original:  v.OriginalURL,
			Shortened: host + "/" + shortUrl,
		}
		urlsToSave = append(urlsToSave, url)
	}

	savedURLs, err := u.repository.InsertBatch(ctx, urlsToSave)
	if err != nil {
		return nil, fmt.Errorf("caller: %s %w", utils.Caller(), err)
	}
	return savedURLs, nil
}
