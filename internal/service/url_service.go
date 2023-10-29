package service

import (
	"fmt"

	"github.com/msmkdenis/yap-shortener/internal/apperrors"
	"github.com/msmkdenis/yap-shortener/internal/model"
	"github.com/msmkdenis/yap-shortener/internal/utils"
	"go.uber.org/zap"
)

type URLService interface {
	Add(s string, host string) (*model.URL, error)
	GetAll() ([]string, error)
	DeleteAll()
	GetByyID(key string) (string, error)
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

func (u *URLUseCase) Add(s, host string) (*model.URL, error) {
	urlKey := utils.GenerateMD5Hash(s)
	url := &model.URL{
		ID:        urlKey,
		Original:  s,
		Shortened: host + "/" + urlKey,
	}

	if savedURL, err := u.repository.SelectByID(url.ID); err == nil {
		return savedURL, nil
	}

	return u.repository.Insert(*url)
}

func (u *URLUseCase) GetAll() ([]string, error) {
	urls, err := u.repository.SelectAll()
	if err != nil {
		return nil, err
	}

	originalURLs := []string{}
	for _, url := range urls {
		originalURLs = append(originalURLs, url.Original)
	}

	return originalURLs, nil
}

func (u *URLUseCase) DeleteAll() {
	u.repository.DeleteAll()
}

func (u *URLUseCase) GetByyID(key string) (string, error) {
	url, err := u.repository.SelectByID(key)
	if err != nil {
		u.logger.Debug(fmt.Sprintf("url with id %s not found", key), zap.Error(err))
		return "", apperrors.ErrorURLNotFound
	}

	return url.Original, nil
}
