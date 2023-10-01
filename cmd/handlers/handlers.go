package handlers

import (
	"fmt"
	"github.com/msmkdenis/yap-shortener/cmd/storage"
	"github.com/msmkdenis/yap-shortener/cmd/utils"
	"io"
	"net/http"
	"strings"
)

func URLHandler(response http.ResponseWriter, request *http.Request) {
	switch request.Method {

	case http.MethodPost:
		body, err := io.ReadAll(request.Body)

		if err != nil {
			http.Error(response, "Unknown Error", http.StatusBadRequest)
		}

		urlKey := utils.GenerateUniqueURLKey()
		storage.Storage[urlKey] = string(body)

		response.WriteHeader(http.StatusCreated)
		response.Write([]byte("http://" + request.Host + "/" + urlKey))

	case http.MethodGet:
		id := (strings.Split(request.URL.Path, "/"))[1]
		url, ok := storage.Storage[id]
		fmt.Println(id)
		fmt.Println(url)
		fmt.Println(storage.Storage)

		if !ok {
			http.Error(response, "Not found", http.StatusBadRequest)
		}

		response.Header().Set("Location", url)
		response.WriteHeader(http.StatusTemporaryRedirect)
	}
}
