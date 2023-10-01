package storage

var GlobalRepository Repository

type URL struct {
	ID        string
	Original  string
	Shortened string
}

type Repository interface {
	Add(u string, host string) URL
	GetByID(key string) (url URL, err error)
}
