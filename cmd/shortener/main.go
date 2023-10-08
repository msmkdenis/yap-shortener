package main

import (
	"github.com/labstack/echo/v4"
	"github.com/msmkdenis/yap-shortener/cmd/handlers"
	"github.com/msmkdenis/yap-shortener/cmd/storage"
)

func main() {
	storage.GlobalRepository = storage.NewMemoryRepository()

	e := echo.New()
	e.POST("/", handlers.PostURL)
	e.GET("/*", handlers.GetURL)
	e.DELETE("/", handlers.DeleteAll)
	e.GET("/", handlers.GetAll)

	e.Logger.Fatal(e.Start(":8080"))
}
