package service

import (
	"fmt"
	"github.com/msmkdenis/yap-shortener/internal/model"
)

type URLService interface {
	Add(s string, host string) model.URL
	GetAll() []string
	DeleteAll()
	GetByyID(key string) (url model.URL, err error)
}

type URLUseCase struct {
	repository model.URLRepository
}

func NewURLService(repository model.URLRepository) *URLUseCase {
	return &URLUseCase{
		repository: repository,
	}
}

func (u *URLUseCase) Add(s string, host string) model.URL {
	url := u.repository.Insert(s, host)
	return url
}

func (u *URLUseCase) GetAll() []string {
	urls := u.repository.SelectAll()
	return urls
}

func (u *URLUseCase) DeleteAll() {
	u.repository.DeleteAll()
}

func (u *URLUseCase) GetByyID(key string) (url model.URL, err error) {
	url, err = u.repository.SelectByID(key)
	if err != nil {
		return url, fmt.Errorf("URL with id = %s not found", key)
	}
	return url, nil
}
