package main

import (
	"github.com/msmkdenis/yap-shortener/cmd/handlers"
	"github.com/msmkdenis/yap-shortener/cmd/storage"
	"net/http"
)

func main() {
	storage.URLRepository = storage.NewURLRepository()
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, handlers.URLHandler)

	err := http.ListenAndServe(`:8080`, mux)

	if err != nil {
		panic(err)
	}
}
