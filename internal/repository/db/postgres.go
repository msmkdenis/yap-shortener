package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PostgresPool struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewPostgresPool(connection string, logger *zap.Logger) *PostgresPool {
	dbPool, err := pgxpool.New(context.Background(), connection)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Unable to connect to database with connection %s", connection), zap.Error(err))
	}

	logger.Info(fmt.Sprintf("Connected to database with connection %s", connection))

	err = dbPool.Ping(context.Background())

	if err == nil {
		logger.Info(fmt.Sprintf("Pinged to database %s", dbPool.Config().ConnConfig.Database))
	}

	return &PostgresPool{
		db:     dbPool,
		logger: logger,
	}
}
