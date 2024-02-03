package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/msmkdenis/yap-shortener/pkg/apperr"
)

// PostgresPool represents PostgreSQL connection pool.
type PostgresPool struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

// NewPostgresPool returns a new instance of PostgresPool with pool of connections.
func NewPostgresPool(connection string, logger *zap.Logger) (*PostgresPool, error) {
	dbPool, err := pgxpool.New(context.Background(), connection)
	if err != nil {
		return nil, apperr.NewValueError(fmt.Sprintf("Unable to connect to database with connection %s", connection), apperr.Caller(), err)
	}

	logger.Info(fmt.Sprintf("Connected to database with connection %s", connection))

	err = dbPool.Ping(context.Background())
	if err != nil {
		return nil, apperr.NewValueError("Unable to ping database", apperr.Caller(), err)
	}
	logger.Info(fmt.Sprintf("Pinged to database %s", dbPool.Config().ConnConfig.Database))

	return &PostgresPool{
		db:     dbPool,
		logger: logger,
	}, nil
}
