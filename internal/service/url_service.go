package service

import (
	"fmt"
	"github.com/labstack/echo/v4"

	"github.com/msmkdenis/yap-shortener/internal/model"
	"github.com/msmkdenis/yap-shortener/internal/utils"
	"go.uber.org/zap"
)

type URLService interface {
	Add(c echo.Context, s string, host string) (*model.URL, error)
	GetAll(c echo.Context) ([]string, error)
	DeleteAll(c echo.Context) error
	GetByyID(c echo.Context, key string) (string, error)
	Ping(c echo.Context) error
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

func (u *URLUseCase) Add(c echo.Context, s, host string) (*model.URL, error) {
	urlKey := utils.GenerateMD5Hash(s)
	url := &model.URL{
		ID:        urlKey,
		Original:  s,
		Shortened: host + "/" + urlKey,
	}

	existingURL, err := u.repository.SelectByID(c, url.ID)
	if err == nil {
		return existingURL, nil
	}

	savedURL, err := u.repository.Insert(c, *url)
	if err != nil {
		return nil, fmt.Errorf("caller: %s %w", utils.Caller(), err)
	}

	return savedURL, nil
}

func (u *URLUseCase) GetAll(c echo.Context) ([]string, error) {
	urls, err := u.repository.SelectAll(c)
	if err != nil {
		return nil, fmt.Errorf("caller: %s %w", utils.Caller(), err)
	}

	originalURLs := []string{}
	for _, url := range urls {
		originalURLs = append(originalURLs, url.Original)
	}

	return originalURLs, nil
}

func (u *URLUseCase) DeleteAll(c echo.Context) error {
	if err := u.repository.DeleteAll(c); err != nil {
		return fmt.Errorf("caller: %s %w", utils.Caller(), err)
	}
	return u.repository.DeleteAll(c)
}

func (u *URLUseCase) GetByyID(c echo.Context, key string) (string, error) {
	url, err := u.repository.SelectByID(c, key)
	if err != nil {
		return "", fmt.Errorf("caller: %s %w", utils.Caller(), err)
	}

	return url.Original, nil
}

func (u *URLUseCase) Ping(c echo.Context) error {
	err := u.repository.Ping(c)
	return err
}
