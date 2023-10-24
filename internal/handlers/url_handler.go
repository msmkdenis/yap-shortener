package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/msmkdenis/yap-shortener/internal/middleware"
	"github.com/msmkdenis/yap-shortener/internal/service"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type URLHandler struct {
	urlService       service.URLService
	urlPrefix        string
	urlHandlerLogger *zap.Logger
}

type URLResponseType struct {
	Result string `json:"result,omitempty"`
}

type URLRequestType struct {
	URL string `json:"url,omitempty"`
}

func New(e *echo.Echo, service service.URLService, urlPrefix string, logger *zap.Logger) *URLHandler {
	handler := &URLHandler{
		urlService:       service,
		urlPrefix:        urlPrefix,
		urlHandlerLogger: logger,
	}

	requestLogger := middleware.InitRequestLogger(logger)

	e.Use(requestLogger.RequestLogger())
	e.Use(middleware.Compress())
	e.Use(middleware.Decompress())

	e.POST("/api/shorten", handler.PostShorten)
	e.POST("/", handler.PostURL)

	e.GET("/*", handler.GetURL)
	e.GET("/", handler.GetAll)

	e.DELETE("/", handler.DeleteAll)

	return handler
}

func (h *URLHandler) PostShorten(c echo.Context) error {

	header := c.Request().Header.Get("Content-Type")
	if header != "application/json" {
		msg := "Content-Type header is not application/json"
		h.urlHandlerLogger.Info("StatusUnsupportedMediaType: " + msg)
		return c.String(http.StatusUnsupportedMediaType, msg)
	}

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		h.urlHandlerLogger.Info("StatusBadRequest: Unknown error, unable to read request")
		return c.String(http.StatusBadRequest, "Error: Unknown error, unable to read request")
	}

	if len(string(body)) == 0 {
		h.urlHandlerLogger.Info("StatusBadRequest: Unable to handle empty body")
		return c.String(http.StatusBadRequest, "Error: Unable to handle empty body")
	}

	var urlRequest URLRequestType
	err = json.Unmarshal(body, &urlRequest)

	if len(urlRequest.URL) == 0 {
		h.urlHandlerLogger.Info("StatusBadRequest: Unable to handle empty body")
		return c.String(http.StatusBadRequest, "Error: Unable to handle empty body")
	}

	if err != nil {
		h.urlHandlerLogger.Info("StatusBadRequest: Unable to unmarshall request")
		return c.String(http.StatusBadRequest, "Error: Unknown error, unable to read request")
	}

	url, _ := h.urlService.Add(urlRequest.URL, h.urlPrefix)

	response := &URLResponseType{
		Result: url.Shortened,
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *URLHandler) PostURL(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)

	if err != nil {
		return c.String(http.StatusBadRequest, "Error: Unknown error, unable to read request")
	}

	if len(string(body)) == 0 {
		return c.String(http.StatusBadRequest, "Error: Unable to handle empty body")
	}

	url, _ := h.urlService.Add(string(body), h.urlPrefix)

	c.Response().WriteHeader(http.StatusCreated)

	return c.String(http.StatusCreated, url.Shortened)
}

func (h *URLHandler) DeleteAll(c echo.Context) error {
	h.urlService.DeleteAll()
	return c.String(http.StatusOK, "All data deleted")
}

func (h *URLHandler) GetAll(c echo.Context) error {
	urls := h.urlService.GetAll()
	return c.String(http.StatusOK, strings.Join(urls, ", "))
}

func (h *URLHandler) GetURL(c echo.Context) error {
	id := (strings.Split(c.Request().URL.Path, "/"))[1]

	if len(id) == 0 {
		return c.String(http.StatusBadRequest, "Error: Unable to handle empty request")
	}

	url, err := h.urlService.GetByyID(id)

	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Error: Not found with id %s", id))
	}

	c.Response().Header().Set("Location", url.Original)
	return c.String(http.StatusTemporaryRedirect, "")
}
