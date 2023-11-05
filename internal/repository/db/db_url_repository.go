package db

import (
	"context"

	"github.com/labstack/echo/v4"

	"go.uber.org/zap"

	"github.com/msmkdenis/yap-shortener/internal/model"
)

type PostgresURLRepository struct {
	PostgresPool *PostgresPool
	logger       *zap.Logger
}

func NewPostgresURLRepository(postgresPool *PostgresPool, logger *zap.Logger) *PostgresURLRepository {
	return &PostgresURLRepository{
		PostgresPool: postgresPool,
		logger:       logger,
	}
}

func (r *PostgresURLRepository) Ping(c echo.Context) error {
	err := r.PostgresPool.db.Ping(context.Background())
	return err
}

func (r *PostgresURLRepository) Insert(c echo.Context, url model.URL) (*model.URL, error) {
	var savedURL model.URL
	err := r.PostgresPool.db.QueryRow(context.Background(),
		`
		insert into url_shortener.url (id, original_url, short_url) 
		values ($1, $2, $3) 
		on conflict (id) do update 
		set original_url = $2, short_url = $3 
		returning id, original_url, short_url
		`,
		url.ID, url.Original, url.Shortened).Scan(&savedURL.ID, &savedURL.Original, &savedURL.Shortened)

	if err != nil {
		r.logger.Error("Query failed", zap.Error(err))
		return nil, err
	}

	return &savedURL, nil
}

func (r *PostgresURLRepository) SelectByID(c echo.Context, key string) (*model.URL, error) {
	var url model.URL
	err := r.PostgresPool.db.QueryRow(context.Background(),
		`
		select id, original_url, short_url
		from url_shortener.url
		where id = $1
		`,
		key).Scan(&url.ID, &url.Original, &url.Shortened)

	if err != nil {
		r.logger.Error("Query failed", zap.Error(err))
		return nil, err
	}

	return &url, nil
}

func (r *PostgresURLRepository) SelectAll(c echo.Context) ([]model.URL, error) {
	var urls []model.URL
	rows, err := r.PostgresPool.db.Query(context.Background(),
		`
		select id, original_url, short_url
		from url_shortener.url
		`)
	if err != nil {
		r.logger.Error("Query failed", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var url model.URL
		err = rows.Scan(&url.ID, &url.Original, &url.Shortened)
		if err != nil {
			r.logger.Error("Query failed", zap.Error(err))
			return nil, err
		}
		urls = append(urls, url)
	}
	if err != nil {
		r.logger.Error("Query failed", zap.Error(err))
		return nil, err
	}
	return urls, nil
}

func (r *PostgresURLRepository) DeleteAll(c echo.Context) error {
	_, err := r.PostgresPool.db.Exec(context.Background(),
	`
	delete from url_shortener.url
	`)

	if err != nil {
		r.logger.Error("Query failed", zap.Error(err))
		return err
	}

	return nil
}
