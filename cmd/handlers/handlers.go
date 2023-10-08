package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/msmkdenis/yap-shortener/cmd/storage"
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

func DeleteAll(c echo.Context) error {
	storage.GlobalRepository.DeleteAll()
	return c.String(http.StatusOK, "All data deleted")
}

func GetAll(c echo.Context) error {
	urls := storage.GlobalRepository.GetAll()
	return c.String(http.StatusOK, strings.Join(urls, " "))
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
