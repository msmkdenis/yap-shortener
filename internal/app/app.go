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
	repository := initRepository(&cfg, logger)
	urlService := service.NewURLService(repository, logger)

	e := echo.New()
	handlers.NewURLHandler(e, urlService, cfg.URLPrefix, logger)
	e.Logger.Fatal(e.Start(cfg.URLServer))
}

func initRepository(cfg *config.Config, logger *zap.Logger) model.URLRepository {
	switch cfg.RepositoryType {
	case config.DataBaseRepository:
		postgresPool, err := db.NewPostgresPool(cfg.DataBaseDSN, logger)
		if err != nil {
			logger.Fatal("Unable to connect to database", zap.Error(err))
		}

		migrations, err := db.NewMigrations(cfg.DataBaseDSN, logger)
		if err != nil {
			logger.Fatal("Unable to create migrations", zap.Error(err))
		}

		err = migrations.MigrateUp()
		if err != nil {
			logger.Fatal("Unable to up migrations", zap.Error(err))
		}

		logger.Info("Connected to database", zap.String("DSN", cfg.DataBaseDSN))
		return db.NewPostgresURLRepository(postgresPool, logger)

	case config.FileRepository:
		repository, err := file.NewFileURLRepository(cfg.FileStoragePath, logger)
		if err != nil {
			logger.Fatal("Unable to create file repository", zap.Error(err))
		}

		logger.Info("Connected/created file", zap.String("FilePath", cfg.FileStoragePath))
		return repository

	default:
		logger.Info("Using memory storage")
		return memory.NewURLRepository(logger)
	}
}
