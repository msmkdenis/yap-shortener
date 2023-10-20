package model

type URL struct {
	ID        string
	Original  string
	Shortened string
}

type URLRepository interface {
	Insert(u string, host string) URL
	SelectByID(key string) (url URL, err error)
	SelectAll() []string
	DeleteAll()
}
