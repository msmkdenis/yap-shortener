package app

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/msmkdenis/yap-shortener/internal/config"
	"github.com/msmkdenis/yap-shortener/internal/handlers"

	"github.com/msmkdenis/yap-shortener/internal/repository/db"
	//"github.com/msmkdenis/yap-shortener/internal/repository/file"
	"github.com/msmkdenis/yap-shortener/internal/service"
	"go.uber.org/zap"
)

func URLShortenerRun() {
	cfg := *config.NewConfig()
	logger, _ := zap.NewProduction()
	//memoryRepository := memory.NewURLRepository(logger)
	//fileRepository := file.NewFileURLRepository(cfg.FileStoragePath, logger)
	postgresPool := db.NewPostgresPool(cfg.DataBaseDSN, logger)
	migrations := db.NewMigrations(cfg.DataBaseDSN, logger)
	logger.Info(fmt.Sprintf("Connecting to database with connection %s", cfg.DataBaseDSN))
	migrations.MigrateUp()
	dbRepository := db.NewPostgresURLRepository(postgresPool, logger)
	//urlService := service.NewURLService(fileRepository, logger)
	urlService := service.NewURLService(dbRepository, logger)
	

	e := echo.New()

	handlers.New(e, urlService, cfg.URLPrefix, logger)

	e.Logger.Fatal(e.Start(cfg.URLServer))
}
