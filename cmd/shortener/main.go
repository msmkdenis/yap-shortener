package main

import (
	"github.com/labstack/echo/v4"
	"github.com/msmkdenis/yap-shortener/internal/config"
	"github.com/msmkdenis/yap-shortener/internal/handlers"
	storage2 "github.com/msmkdenis/yap-shortener/internal/storage"
)

func main() {

	config.AppConfig = config.InitConfig()

	storage2.GlobalRepository = storage2.NewMemoryRepository()

	e := echo.New()
	e.POST("/", handlers.PostURL)
	e.GET("/*", handlers.GetURL)
	e.DELETE("/", handlers.DeleteAll)
	e.GET("/", handlers.GetAll)

	e.Logger.Fatal(e.Start(config.AppConfig.URLServer))
}
