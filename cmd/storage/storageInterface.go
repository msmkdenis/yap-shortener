package storage

type UrlStorage interface {
	Add(u string, host string) URL
	GetById(key string) (url URL, err error)
}
