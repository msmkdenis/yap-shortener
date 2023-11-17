package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/msmkdenis/yap-shortener/internal/apperrors"
	"github.com/msmkdenis/yap-shortener/internal/utils"
	"go.uber.org/zap"
)

type PostgresPool struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewPostgresPool(connection string, logger *zap.Logger) (*PostgresPool, error) {
	dbPool, err := pgxpool.New(context.Background(), connection)
	if err != nil {
		return nil, apperrors.NewValueError(fmt.Sprintf("Unable to connect to database with connection %s", connection), utils.Caller(), err)
	}

	logger.Info(fmt.Sprintf("Connected to database with connection %s", connection))

	err = dbPool.Ping(context.Background())
	if err != nil {
		return nil, apperrors.NewValueError("Unable to ping database", utils.Caller(), err)
	}
	logger.Info(fmt.Sprintf("Pinged to database %s", dbPool.Config().ConnConfig.Database))

	return &PostgresPool{
		db:     dbPool,
		logger: logger,
	}, nil
}
