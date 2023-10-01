package main

import (
	"io"
	"math/rand"
	"net/http"
	"strings"
)

var storage map[string]string

func generateUniqueUrlKey() string {
	runes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	urlKey := make([]rune, 8)
	for i := range urlKey {
		urlKey[i] = runes[rand.Intn(len(runes))]
	}

	// Проверка на уникальность через рекурсию, пока не создастся уникальный ключ
	_, ok := storage[string(urlKey)]
	if ok {
		return generateUniqueUrlKey()
	}

	return string(urlKey)
}

func URLHandler(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		body, err := io.ReadAll(request.Body)

		if err != nil {
			http.Error(response, "Unknown Error", http.StatusBadRequest)
		}

		urlKey := generateUniqueUrlKey()
		storage[urlKey] = string(body)

		response.WriteHeader(http.StatusCreated)
		response.Write([]byte("http://" + request.Host + "/" + urlKey))
	case http.MethodGet:
		id := (strings.Split(request.URL.Path, "/"))[1]
		url, ok := storage[id]
		if !ok {
			http.Error(response, "Not found", http.StatusBadRequest)
		}
		response.Header().Set("Location", url)
		response.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func main() {
	storage = make(map[string]string)
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, URLHandler)

	err := http.ListenAndServe(`:8080`, mux)

	if err != nil {
		panic(err)
	}
}
