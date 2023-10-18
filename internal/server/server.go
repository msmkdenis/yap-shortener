package server

import (
	"github.com/labstack/echo/v4"
	"github.com/msmkdenis/yap-shortener/internal/handlers"
)

func InitServer(URLServer string) {
	e := echo.New()
	e.POST("/", handlers.PostURL)
	e.GET("/*", handlers.GetURL)
	e.DELETE("/", handlers.DeleteAll)
	e.GET("/", handlers.GetAll)

	e.Logger.Fatal(e.Start(URLServer))
}
