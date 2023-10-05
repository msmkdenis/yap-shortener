package main

import (
	"github.com/labstack/echo/v4"
	"github.com/msmkdenis/yap-shortener/cmd/handlers"
	"github.com/msmkdenis/yap-shortener/cmd/storage"
)

func main() {
	storage.GlobalRepository = storage.NewMemoryRepository()
	/*	mux := http.NewServeMux()
		mux.HandleFunc(`/`, handlers.URLHandler)

		err := http.ListenAndServe(`:8080`, mux)

		if err != nil {
			panic(err)
		}*/

	e := echo.New()
	e.POST("/", handlers.PostURL)
	e.GET("/*", handlers.GetURL)

	e.Logger.Fatal(e.Start(":8080"))
}
