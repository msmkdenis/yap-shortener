package service

import (
	"fmt"

	"github.com/msmkdenis/yap-shortener/internal/model"
	"github.com/msmkdenis/yap-shortener/internal/utils"
	"go.uber.org/zap"
)

type URLService interface {
	Add(s string, host string) (*model.URL, error)
	GetAll() ([]string, error)
	DeleteAll() error
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

	existingURL, err := u.repository.SelectByID(url.ID)
	if err == nil {
		return existingURL, nil
	}

	savedURL, err := u.repository.Insert(*url)
	if err != nil {
		return nil, fmt.Errorf("caller: %s %w", utils.Caller(), err)
	}

	return savedURL, nil
}

func (u *URLUseCase) GetAll() ([]string, error) {
	urls, err := u.repository.SelectAll()
	if err != nil {
		return nil, fmt.Errorf("caller: %s %w", utils.Caller(), err)
	}

	originalURLs := []string{}
	for _, url := range urls {
		originalURLs = append(originalURLs, url.Original)
	}

	return originalURLs, nil
}

func (u *URLUseCase) DeleteAll() error {
	if err := u.repository.DeleteAll(); err != nil { 
		return fmt.Errorf("caller: %s %w", utils.Caller(), err) 
	} 
	return u.repository.DeleteAll()
}

func (u *URLUseCase) GetByyID(key string) (string, error) {
	url, err := u.repository.SelectByID(key)
	if err != nil {
		return "", fmt.Errorf("caller: %s %w", utils.Caller(), err)
	}

	return url.Original, nil
}
