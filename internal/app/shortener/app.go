package shortener

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/msmkdenis/yap-shortener/pkg/echopprof"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"

	"github.com/msmkdenis/yap-shortener/internal/config"
	"github.com/msmkdenis/yap-shortener/internal/handlers"
	"github.com/msmkdenis/yap-shortener/internal/repository/db"
	"github.com/msmkdenis/yap-shortener/internal/repository/file"
	"github.com/msmkdenis/yap-shortener/internal/repository/memory"
	"github.com/msmkdenis/yap-shortener/internal/service"
	"github.com/msmkdenis/yap-shortener/internal/utils"
)

func URLShortenerRun() {
	cfg := *config.NewConfig()
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Unable to initialize zap logger", zap.Error(err))
	}
	jwtManager := utils.InitJWTManager(cfg.TokenName, cfg.SecretKey, logger)
	repository := initRepository(&cfg, logger)
	urlService := service.NewURLService(repository, logger)

	e := echo.New()
	echopprof.Wrap(e)
	handlers.NewURLHandler(e, urlService, cfg.URLPrefix, jwtManager, logger)

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-quit

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 5*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		if errShutdown := e.Shutdown(shutdownCtx); errShutdown != nil {
			e.Logger.Fatal(errShutdown)
		}
		serverStopCtx()
	}()

	errStart := e.Start(cfg.URLServer)
	if errStart != nil && !errors.Is(errStart, http.ErrServerClosed) {
		log.Fatal(err)
	}

	<-serverCtx.Done()
}

func initRepository(cfg *config.Config, logger *zap.Logger) service.URLRepository {
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
