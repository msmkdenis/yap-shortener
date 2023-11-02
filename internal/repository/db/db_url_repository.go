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
	return nil, nil
}

func (r *PostgresURLRepository) SelectByID(c echo.Context, key string) (*model.URL, error) {
	return nil, nil
}

func (r *PostgresURLRepository) SelectAll(c echo.Context) ([]model.URL, error) {
	return nil, nil
}

func (r *PostgresURLRepository) DeleteAll(c echo.Context) error {
	return nil
}
