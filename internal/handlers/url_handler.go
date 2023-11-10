package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/msmkdenis/yap-shortener/internal/apperrors"
	"github.com/msmkdenis/yap-shortener/internal/handlers/dto"
	"github.com/msmkdenis/yap-shortener/internal/middleware"
	"github.com/msmkdenis/yap-shortener/internal/service"
	"github.com/msmkdenis/yap-shortener/internal/utils"
	"go.uber.org/zap"

	"github.com/labstack/echo/v4"
)

type URLHandler struct {
	urlService service.URLService
	urlPrefix  string
	logger     *zap.Logger
}

func New(e *echo.Echo, service service.URLService, urlPrefix string, logger *zap.Logger) *URLHandler {
	handler := &URLHandler{
		urlService: service,
		urlPrefix:  urlPrefix,
		logger:     logger,
	}

	requestLogger := middleware.InitRequestLogger(logger)

	e.Use(requestLogger.RequestLogger())
	e.Use(middleware.Compress())
	e.Use(middleware.Decompress())

	e.POST("/api/shorten", handler.PostShorten)
	e.POST("/", handler.PostURL)
	e.POST("/api/shorten/batch", handler.PostBatch)

	e.GET("/*", handler.GetURL)
	e.GET("/", handler.GetAll)
	e.GET("/ping", handler.GetPing)

	e.DELETE("/", handler.DeleteAll)

	return handler
}

func (h *URLHandler) PostBatch(c echo.Context) error {
	header := c.Request().Header.Get("Content-Type")
	if header != "application/json" {
		msg := "Content-Type header is not application/json"
		h.logger.Error("StatusUnsupportedMediaType: " + msg)
		return c.String(http.StatusUnsupportedMediaType, msg)
	}

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		h.logger.Error("StatusBadRequest: unknown error", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), err)))
		return c.String(http.StatusBadRequest, fmt.Sprintf("Error: Unknown error, unable to read request %s", err))
	}

	var urlBatchRequest []dto.URLBatchRequestType
	err = json.Unmarshal([]byte(body), &urlBatchRequest)
	if err != nil {
		h.logger.Error("StatusBadRequest: unknown error", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), err)))
		return c.String(http.StatusBadRequest, "Error: Unknown error, unable to read request")
	}

	if err := h.checkRequest(string(body)); err != nil {
		h.logger.Error("StatusBadRequest: unable to handle empty request", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), err)))
		return c.String(http.StatusBadRequest, "Error: Unable to handle empty request")
	}

	savedURLs, err := h.urlService.AddAll(c, urlBatchRequest, h.urlPrefix)
	if err != nil {
		h.logger.Error("StatusInternalServerError: unknown error", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), err)))
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unknown error: %s", err))
	}

	return c.JSON(http.StatusCreated, savedURLs)
}

func (h *URLHandler) PostShorten(c echo.Context) error {
	header := c.Request().Header.Get("Content-Type")
	if header != "application/json" {
		msg := "Content-Type header is not application/json"
		h.logger.Error("StatusUnsupportedMediaType: " + msg)
		return c.String(http.StatusUnsupportedMediaType, msg)
	}

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		h.logger.Error("StatusBadRequest: unknown error", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), err)))
		return c.String(http.StatusBadRequest, fmt.Sprintf("Error: Unknown error, unable to read request %s", err))
	}

	var urlRequest dto.URLRequestType
	err = json.Unmarshal(body, &urlRequest)
	if err != nil {
		h.logger.Error("StatusBadRequest: unknown error", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), err)))
		return c.String(http.StatusBadRequest, "Error: Unknown error, unable to read request")
	}

	if err := h.checkRequest(urlRequest.URL); err != nil {
		h.logger.Error("StatusBadRequest: unable to handle empty request", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), err)))
		return c.String(http.StatusBadRequest, "Error: Unable to handle empty request")
	}

	url, err := h.urlService.Add(c, urlRequest.URL, h.urlPrefix)
	if err != nil {
		h.logger.Error("StatusBadRequest: unknown error", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), err)))
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unknown error: %s", err))
	}

	response := &dto.URLResponseType{
		Result: url.Shortened,
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *URLHandler) PostURL(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		h.logger.Error("StatusBadRequest: unknown error", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), err)))
		return c.String(http.StatusBadRequest, "Error: Unknown error, unable to read request")
	}

	if err := h.checkRequest(string(body)); err != nil {
		h.logger.Error("StatusBadRequest: unable to handle empty request", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), err)))
		return c.String(http.StatusBadRequest, "Error: Unable to handle empty request")
	}

	url, err := h.urlService.Add(c, string(body), h.urlPrefix)
	if err != nil {
		h.logger.Error("StatusInternalServerError: Unknown error:", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), err)))
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unknown error: %s", err))
	}

	c.Response().WriteHeader(http.StatusCreated)
	return c.String(http.StatusCreated, url.Shortened)
}

func (h *URLHandler) DeleteAll(c echo.Context) error {
	if err := h.urlService.DeleteAll(c); err != nil {
		h.logger.Error("StatusInternalServerError: Unknown error:", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), err)))
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unknown error: %s", err))
	}

	return c.String(http.StatusOK, "All data deleted")
}

func (h *URLHandler) GetAll(c echo.Context) error {
	urls, err := h.urlService.GetAll(c)
	if err != nil {
		h.logger.Error("StatusInternalServerError: Unknown error:", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), err)))
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unknown error: %s", err))
	}

	return c.String(http.StatusOK, strings.Join(urls, ", "))
}

func (h *URLHandler) GetURL(c echo.Context) error {
	id := (strings.Split(c.Request().URL.Path, "/"))[1]

	if err := h.checkRequest(id); err != nil {
		h.logger.Error("StatusBadRequest: Unable to handle empty request", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), err)))
		return c.String(http.StatusBadRequest, "Error: Unable to handle empty request")
	}

	originalURL, err := h.urlService.GetByyID(c, id)

	var message string
	var status int

	switch {
	case errors.Is(err, apperrors.ErrorURLNotFound):
		h.logger.Info("StatusBadRequest: url not found", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), err)))
		status = http.StatusBadRequest
		message = fmt.Sprintf("URL with id %s not found", id)

	case err != nil:
		h.logger.Error("InternalServerError", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), err)))
		status = http.StatusInternalServerError
		message = fmt.Sprintf("Unknown error: %s", err)

	default:
		c.Response().Header().Set("Location", originalURL)
		status = http.StatusTemporaryRedirect
		message = ""
	}
	fmt.Println("++++++++++++++++++++++++++++")
	fmt.Println(c.Response().Header().Get("Location"))
	fmt.Println(status)

	return c.String(status, message)
}

func (h *URLHandler) checkRequest(s string) error {
	if len(s) == 0 {
		return apperrors.NewValueError("Unable to handle empty request", utils.Caller(), apperrors.ErrorEmptyRequest)
	}

	return nil
}

func (h *URLHandler) GetPing(c echo.Context) error {
	status := http.StatusOK
	err := h.urlService.Ping(c)
	if err != nil {
		status = http.StatusInternalServerError
	}

	return c.NoContent(status)
}
