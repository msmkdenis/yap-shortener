package service

import (
	"fmt"
	"github.com/msmkdenis/yap-shortener/internal/model"
	"github.com/msmkdenis/yap-shortener/internal/utils"
	"go.uber.org/zap"
)

type URLService interface {
	Add(s string, host string) (*model.URL, error)
	GetAll() []string
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

func (u *URLUseCase) Add(s string, host string) (*model.URL, error) {

	urlKey := utils.GenerateMD5Hash(s)

	var url = &model.URL{
		ID:        urlKey,
		Original:  s,
		Shortened: host + "/" + urlKey,
	}

	savedURL, err := u.repository.Insert(*url)

	if err != nil {
		return nil, err
	}

	return savedURL, nil
}

func (u *URLUseCase) GetAll() []string {
	urls, _ := u.repository.SelectAll()
	var originalURLs []string
	for _, url := range urls {
		originalURLs = append(originalURLs, url.Original)
	}
	return originalURLs
}

func (u *URLUseCase) DeleteAll() {
	u.repository.DeleteAll()
}

func (u *URLUseCase) GetByyID(key string) (string, error) {
	var url = &model.URL{}
	url, err := u.repository.SelectByID(key)
	if err != nil {
		return "", fmt.Errorf("URL with id = %s not found", key)
	}
	return url.Original, nil
}
