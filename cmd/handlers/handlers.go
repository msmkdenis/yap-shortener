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
			return
		}

		if len(string(body)) == 0 {
			http.Error(response, "Cant handle empty body!", http.StatusBadRequest)
			return
		}

		url := storage.GlobalRepository.Add(string(body), request.Host)

		response.WriteHeader(http.StatusCreated)
		response.Write([]byte(url.Shortened))

	case http.MethodGet:
		id := (strings.Split(request.URL.Path, "/"))[1]

		if len(id) == 0 {
			http.Error(response, "Cant handle empty request!", http.StatusBadRequest)
			return
		}

		url, err := storage.GlobalRepository.GetByID(id)

		if err != nil {
			http.Error(response, "Not found", http.StatusBadRequest)
			return
		}

		response.Header().Set("Location", url.Original)
		response.WriteHeader(http.StatusTemporaryRedirect)
	}
}
