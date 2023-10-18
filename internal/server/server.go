package server

import (
	"github.com/labstack/echo/v4"
	"github.com/msmkdenis/yap-shortener/internal/handlers"
	"github.com/msmkdenis/yap-shortener/internal/middleware"
	"go.uber.org/zap"
)

func InitServer(URLServer string) {
	logger, _ := zap.NewProduction()
	requestLogger := middleware.InitRequestLogger(logger)
	e := echo.New()
	e.Use(requestLogger.RequestLogger())
	e.POST("/", handlers.PostURL)
	e.GET("/*", handlers.GetURL)
	e.DELETE("/", handlers.DeleteAll)
	e.GET("/", handlers.GetAll)

	e.Logger.Fatal(e.Start(URLServer))
}
