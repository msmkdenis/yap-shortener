package db

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

	pwd, _ := os.Getwd()
	filedir := filepath.Dir(filepath.Dir(pwd))
	logger.Info(filedir)

	file, err := os.Open(filepath.Join(filedir, "schema.sql"))
	if err != nil {
		logger.Fatal("Unable to read schema.sql file", zap.Error(err))
	}

	data, _ := os.ReadFile(file.Name())

	dbPool.Exec(context.Background(), strings.TrimSpace(string(data)))

	return &PostgresPool{
		db:     dbPool,
		logger: logger,
	}
}
