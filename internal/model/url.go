package model

import (
	"context"
)

type URL struct {
	ID            string `db:"id"`
	Original      string `db:"original_url"`
	Shortened     string `db:"short_url"`
	CorrelationID string `db:"correlation_id"`
}

type URLRepository interface {
	Insert(ctx context.Context, u URL) (*URL, error)
	InsertAllOrUpdate(ctx context.Context, urls []URL) ([]URL, error)
	SelectByID(ctx context.Context, key string) (*URL, error)
	SelectAll(ctx context.Context) ([]URL, error)
	DeleteAll(ctx context.Context) error
	Ping(ctx context.Context) error
}
