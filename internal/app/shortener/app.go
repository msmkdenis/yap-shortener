// Package shortener implements the URL shortener service.
package shortener

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/msmkdenis/yap-shortener/internal/api/grpchandlers"
	"github.com/msmkdenis/yap-shortener/internal/api/httphandlers"
	"github.com/msmkdenis/yap-shortener/internal/config"
	"github.com/msmkdenis/yap-shortener/internal/middleware"
	pb "github.com/msmkdenis/yap-shortener/internal/proto"
	"github.com/msmkdenis/yap-shortener/internal/repository/db"
	"github.com/msmkdenis/yap-shortener/internal/repository/file"
	"github.com/msmkdenis/yap-shortener/internal/repository/memory"
	"github.com/msmkdenis/yap-shortener/internal/service"
	"github.com/msmkdenis/yap-shortener/pkg/echopprof"
	"github.com/msmkdenis/yap-shortener/pkg/jwtgen"
)

// URLShortenerRun runs the URL shortener service. Graceful shutdown is implemented.
//
// It does not take any parameters and does not return any values.
func URLShortenerRun() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Unable to initialize zap logger", zap.Error(err))
	}

	cfg := *config.NewConfig(logger)

	jwtManager := jwtgen.InitJWTManager(cfg.TokenName, cfg.SecretKey, logger)
	jwtCheckerCreator := middleware.InitJWTCheckerCreator(jwtManager, logger)
	jwtAuth := middleware.InitJWTAuth(jwtManager, logger)
	repository := initRepository(&cfg, logger)
	urlService := service.NewURLService(repository, logger)

	e := echo.New()
	echopprof.Wrap(e)
	wgHTTP := &sync.WaitGroup{}
	httphandlers.NewURLShorten(e, urlService, cfg.URLPrefix, cfg.TrustedSubnet, jwtCheckerCreator, jwtAuth, logger, wgHTTP)

	listener, err := net.Listen("tcp", cfg.GRPCServer)
	if err != nil {
		logger.Fatal("Unable to create listener", zap.Error(err))
	}
	serverGrpc := grpc.NewServer(
		grpc.ChainUnaryInterceptor(jwtAuth.GRPCJWTAuth, jwtCheckerCreator.GRPCJWTCheckOrCreate),
	)
	wgGRPC := &sync.WaitGroup{}
	pb.RegisterURLShortenerServer(serverGrpc, grpchandlers.NewURLShorten(urlService, cfg.URLPrefix, cfg.TrustedSubnet, jwtManager, logger, wgGRPC))
	reflection.Register(serverGrpc)

	httpServerCtx, httpServerStopCtx := context.WithCancel(context.Background())
	grpcServerCtx, grpcServerStopCtx := context.WithCancel(context.Background())

	// Канал для сигналов
	quitSignal := make(chan os.Signal, 1)
	signal.Notify(quitSignal, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	quit := make(chan struct{})
	go func() {
		// Получили сигнал
		<-quitSignal
		// Закрыли сигнальный канал
		close(quit)
	}()

	// Запустили сервер gRPC
	go func() {
		logger.Info(fmt.Sprintf("gRPC server starting on port %s", cfg.GRPCServer))
		if errGRPC := serverGrpc.Serve(listener); errGRPC != nil {
			logger.Fatal("Unable to start gRPC server", zap.Error(errGRPC))
		}
	}()

	// Запустили сервер HTTP
	go func() {
		if cfg.EnableHTTPS == "true" {
			e.AutoTLSManager.Cache = autocert.DirCache("cache-dir")
			errStart := e.StartAutoTLS(cfg.URLServer)
			if errStart != nil && !errors.Is(errStart, http.ErrServerClosed) {
				log.Fatal(err)
			}
		} else {
			errStart := e.Start(cfg.URLServer)
			if errStart != nil && !errors.Is(errStart, http.ErrServerClosed) {
				log.Fatal(err)
			}
		}
	}()

	// Graceful shutdown ПКЗС
	go func() {
		// Слушаем сигнальный канал, при закрытии код идет дальше
		<-quit

		// Shutdown signal with grace period of 10 seconds
		shutdownCtx, cancel := context.WithTimeout(httpServerCtx, 10*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		logger.Info("Shutdown signal received, gracefully stopping httphandlers server")
		if errShutdown := e.Shutdown(shutdownCtx); errShutdown != nil {
			e.Logger.Fatal(errShutdown)
		}
		httpServerStopCtx()
	}()

	go func() {
		<-quit

		// Shutdown signal with grace period of 10 seconds
		shutdownCtx, cancel := context.WithTimeout(grpcServerCtx, 10*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				logger.Error("graceful gRPC shutdown timed out.. forcing exit.")
				serverGrpc.Stop()
			}
		}()

		// Trigger graceful shutdown
		logger.Info("Shutdown signal received, gracefully stopping gRPC server")
		serverGrpc.GracefulStop()
		grpcServerStopCtx()
	}()

	wgHTTP.Wait()
	wgGRPC.Wait()
	<-httpServerCtx.Done()
	<-grpcServerCtx.Done()
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
