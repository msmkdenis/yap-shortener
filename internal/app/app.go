package app

import (
	"github.com/labstack/echo/v4"
	"github.com/msmkdenis/yap-shortener/internal/config"
	"github.com/msmkdenis/yap-shortener/internal/handlers"
	"github.com/msmkdenis/yap-shortener/internal/model"

	"github.com/msmkdenis/yap-shortener/internal/repository/db"
	"github.com/msmkdenis/yap-shortener/internal/repository/file"
	"github.com/msmkdenis/yap-shortener/internal/repository/memory"
	"github.com/msmkdenis/yap-shortener/internal/service"
	"go.uber.org/zap"
)

func URLShortenerRun() {
	cfg := *config.NewConfig()
	logger, _ := zap.NewProduction()

	var repository model.URLRepository
	if cfg.RepositoryType.DataBaseRepository {
		postgresPool := db.NewPostgresPool(cfg.DataBaseDSN, logger)
		migrations := db.NewMigrations(cfg.DataBaseDSN, logger)
		migrations.MigrateUp()
		repository = db.NewPostgresURLRepository(postgresPool, logger)
		logger.Info("Connected to database", zap.String("DSN", cfg.DataBaseDSN))
	} else if cfg.RepositoryType.FileRepository {
		repository = file.NewFileURLRepository(cfg.FileStoragePath, logger)
		logger.Info("Connected/created file", zap.String("FilePath", cfg.FileStoragePath))
	} else {
		repository = memory.NewURLRepository(logger)
		logger.Info("Using memory storage")
	}
	urlService := service.NewURLService(repository, logger)
	

	e := echo.New()

	handlers.New(e, urlService, cfg.URLPrefix, logger)

	e.Logger.Fatal(e.Start(cfg.URLServer))
}
