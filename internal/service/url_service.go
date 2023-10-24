package service

import (
	"fmt"
	"github.com/msmkdenis/yap-shortener/internal/model"
	"github.com/msmkdenis/yap-shortener/internal/utils"
)

type URLService interface {
	Add(s string, host string) (*model.URL, error)
	GetAll() []string
	DeleteAll()
	GetByyID(key string) (*model.URL, error)
}

type URLUseCase struct {
	repository model.URLRepository
}

func NewURLService(repository model.URLRepository) *URLUseCase {
	return &URLUseCase{
		repository: repository,
	}
}

func (u *URLUseCase) Add(s string, host string) (*model.URL, error) {

	urlKey := utils.GenerateMD5Hash(s)

	var url = &model.URL{
		ID:        urlKey,
		Original:  s,
		Shortened: host + "/" + urlKey,
	}

	savedUrl, err := u.repository.Insert(*url)

	if err != nil {
		return nil, err
	}

	return savedUrl, nil
}

func (u *URLUseCase) GetAll() []string {
	urls, _ := u.repository.SelectAll()
	return urls
}

func (u *URLUseCase) DeleteAll() {
	u.repository.DeleteAll()
}

func (u *URLUseCase) GetByyID(key string) (*model.URL, error) {
	var url = &model.URL{}
	url, err := u.repository.SelectByID(key)
	if err != nil {
		return url, fmt.Errorf("URL with id = %s not found", key)
	}
	return url, nil
}
