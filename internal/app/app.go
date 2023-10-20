package app

import (
	"github.com/labstack/echo/v4"
	"github.com/msmkdenis/yap-shortener/internal/config"
	"github.com/msmkdenis/yap-shortener/internal/handlers"
	"github.com/msmkdenis/yap-shortener/internal/repository/memory"
	"github.com/msmkdenis/yap-shortener/internal/service"
	"go.uber.org/zap"
)

func URLShortenerRun() {
	cfg := *config.NewConfig()
	urlRepository := memory.NewURLRepository()
	urlService := service.NewURLService(urlRepository)
	logger, _ := zap.NewProduction()

	e := echo.New()

	handlers.New(e, urlService, cfg.URLPrefix, logger)

	e.Logger.Fatal(e.Start(cfg.URLServer))
}
