package model

import "github.com/labstack/echo/v4"

type URL struct {
	ID        string `db:"id"`
	Original  string `db:"original_url"`
	Shortened string `db:"short_url"`
}

type URLRepository interface {
	InsertOrUpdate(c echo.Context, u URL) (*URL, error)
	InsertAllOrUpdate(c echo.Context, urls []URL) ([]URL, error)
	SelectByID(c echo.Context, key string) (*URL, error)
	SelectAll(c echo.Context) ([]URL, error)
	DeleteAll(c echo.Context) error
	Ping(c echo.Context) error
}
