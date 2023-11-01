package model

type URL struct {
	ID        string
	Original  string
	Shortened string
}

type URLRepository interface {
	Insert(u URL) (*URL, error)
	SelectByID(key string) (*URL, error)
	SelectAll() ([]URL, error)
	DeleteAll() error
}
