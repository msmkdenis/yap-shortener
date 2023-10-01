package handlers

import (
	"github.com/msmkdenis/yap-shortener/cmd/storage"
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

		url := storage.URLRepository.Add(string(body), request.Host)

		/*		urlKey := utils.GenerateUniqueURLKey()
				storage.Storage[urlKey] = string(body)

				response.WriteHeader(http.StatusCreated)
				response.Write([]byte("http://" + request.Host + "/" + urlKey))*/

		response.WriteHeader(http.StatusCreated)
		response.Write([]byte(url.Shortened))

	case http.MethodGet:
		id := (strings.Split(request.URL.Path, "/"))[1]

		url, err := storage.URLRepository.GetById(id)

		if err != nil {
			http.Error(response, "Not found", http.StatusBadRequest)
		}

		/*		url, ok := storage.Storage[id]

				if !ok {
					http.Error(response, "Not found", http.StatusBadRequest)
				}*/

		response.Header().Set("Location", url.Original)
		response.WriteHeader(http.StatusTemporaryRedirect)
	}
}
