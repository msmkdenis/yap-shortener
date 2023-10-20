package service

import (
	"fmt"
	"github.com/msmkdenis/yap-shortener/internal/domain"
)

type URLService interface {
	Add(s string, host string) domain.URL
	GetAll() []string
	DeleteAll()
	GetByyID(key string) (url domain.URL, err error)
}

type URLUseCase struct {
	repository domain.URLRepository
}

func NewURLService(repository domain.URLRepository) *URLUseCase {
	return &URLUseCase{
		repository: repository,
	}
}

func (u *URLUseCase) Add(s string, host string) domain.URL {
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

func (u *URLUseCase) GetByyID(key string) (url domain.URL, err error) {
	url, err = u.repository.SelectByID(key)
	if err != nil {
		return url, fmt.Errorf("URL with id = %s not found", key)
	}
	return url, nil
}
