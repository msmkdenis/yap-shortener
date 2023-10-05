package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/msmkdenis/yap-shortener/cmd/storage"
	"io"
	"net/http"
	"strings"
)

func PostURL(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)

	if err != nil {
		return c.String(http.StatusBadRequest, "Error: Unknown error, unable to read request")
	}

	if len(string(body)) == 0 {
		return c.String(http.StatusBadRequest, "Error: Unable to handle empty body")
	}

	url := storage.GlobalRepository.Add(string(body), c.Request().Host)

	c.Response().WriteHeader(http.StatusCreated)

	return c.String(http.StatusCreated, url.Shortened)
}

func GetURL(c echo.Context) error {
	id := (strings.Split(c.Request().URL.Path, "/"))[1]

	if len(id) == 0 {
		return c.String(http.StatusBadRequest, "Error: Unable to handle empty request")
	}

	url, err := storage.GlobalRepository.GetByID(id)

	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Error: Not found with id %s", id))
	}

	c.Response().Header().Set("Location", url.Original)
	return c.String(http.StatusTemporaryRedirect, "")
}

/*func URLHandler(response http.ResponseWriter, request *http.Request) {
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
}*/
