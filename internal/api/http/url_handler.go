// Package http implements the URL shortener service http handlers.
package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"

	"github.com/msmkdenis/yap-shortener/internal/dto"
	"github.com/msmkdenis/yap-shortener/internal/middleware"
	"github.com/msmkdenis/yap-shortener/internal/model"
	urlErr "github.com/msmkdenis/yap-shortener/internal/urlerr"
	"github.com/msmkdenis/yap-shortener/pkg/apperr"
	"github.com/msmkdenis/yap-shortener/pkg/workerpool"
)

// URLHandler represents URL handler struct.
type URLHandler struct {
	urlService    URLService
	urlPrefix     string
	trustedSubnet string
	logger        *zap.Logger
	wg            *sync.WaitGroup
}

// URLService represents URL service interface.
type URLService interface {
	Add(ctx context.Context, s string, host string, userID string) (*model.URL, error)
	AddAll(ctx context.Context, urls []dto.URLBatchRequest, host string, userID string) ([]dto.URLBatchResponse, error)
	GetAll(ctx context.Context) ([]string, error)
	GetAllByUserID(ctx context.Context, userID string) ([]dto.URLBatchResponseByUserID, error)
	DeleteAll(ctx context.Context) error
	DeleteURLByUserID(ctx context.Context, userID string, shortURLs string) error
	GetByyID(ctx context.Context, key string) (string, error)
	GetStats(ctx context.Context) (*dto.URLStats, error)
	Ping(ctx context.Context) error
}

// NewURLHandler creates a new URLHandler instance
//
// Registers the URL shortener service http handlers.
func NewURLHandler(e *echo.Echo, service URLService, urlPrefix string, trustedSubnet string, jwtCheckerCreator *middleware.JWTCheckerCreator, jwtAuth *middleware.JWTAuth, logger *zap.Logger, wg *sync.WaitGroup) *URLHandler {
	handler := &URLHandler{
		urlService:    service,
		urlPrefix:     urlPrefix,
		trustedSubnet: trustedSubnet,
		logger:        logger,
		wg:            wg,
	}

	requestLogger := middleware.InitRequestLogger(logger)

	e.Use(requestLogger.RequestLogger())
	e.Use(middleware.Compress())
	e.Use(middleware.Decompress())

	public := e.Group("/", jwtCheckerCreator.JWTCheckOrCreate())
	public.POST("api/shorten", handler.AddShorten)
	public.POST("", handler.AddURL)
	public.POST("api/shorten/batch", handler.AddBatch)

	public.GET("*", handler.FindURL)
	public.GET("", handler.FindAll)
	public.GET("ping", handler.Ping)

	public.DELETE("", handler.ClearAll)

	protected := e.Group("/api/user", jwtAuth.JWTAuth())
	protected.GET("/urls", handler.FindAllURLByUserID)
	protected.DELETE("/urls", handler.DeleteAllURLsByUserID)

	e.GET("/api/internal/stats", handler.GetStats)

	return handler
}

// FindAllURLByUserID retrieves all URLs for a given user ID.
func (h *URLHandler) FindAllURLByUserID(c echo.Context) error {
	userID, ok := c.Get("userID").(string)
	if !ok {
		h.logger.Error("Internal server error", zap.Error(urlErr.ErrUnableToGetUserIDFromContext))
		return c.NoContent(http.StatusInternalServerError)
	}

	savedURLs, err := h.urlService.GetAllByUserID(c.Request().Context(), userID)
	if err != nil && !errors.Is(err, urlErr.ErrURLNotFound) {
		h.logger.Error("StatusBadRequest: unknown error", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), err)))
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unknown error: %s", err))
	}

	if errors.Is(err, urlErr.ErrURLNotFound) {
		h.logger.Warn("StatusNoContent: urls not found", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), err)))
		return c.NoContent(http.StatusNoContent)
	}

	return c.JSON(http.StatusOK, savedURLs)
}

// GetStats returns URL stats.
func (h *URLHandler) GetStats(c echo.Context) error {
	if h.trustedSubnet == "" {
		return c.NoContent(http.StatusForbidden)
	}

	ip := c.Request().Header.Get("X-Real-IP")
	if ip == "" {
		return c.NoContent(http.StatusForbidden)
	}

	_, ipNet, err := net.ParseCIDR(h.trustedSubnet)
	if err != nil {
		h.logger.Error("StatusInternalServerError: unable to parse CIDR", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), err)))
		return c.NoContent(http.StatusInternalServerError)
	}

	if !ipNet.Contains(net.ParseIP(ip)) {
		return c.NoContent(http.StatusForbidden)
	}

	stats, err := h.urlService.GetStats(c.Request().Context())
	if err != nil {
		h.logger.Error("StatusInternalServerError: unknown error", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), err)))
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, stats)
}

// DeleteAllURLsByUserID deletes all URLs associated with a user ID.
func (h *URLHandler) DeleteAllURLsByUserID(c echo.Context) error {
	header := c.Request().Header.Get("Content-Type")
	if header != "application/json" {
		msg := "Content-Type header is not application/json"
		h.logger.Error("StatusUnsupportedMediaType: " + msg)
		return c.String(http.StatusUnsupportedMediaType, msg)
	}

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		h.logger.Error("StatusBadRequest: unknown error", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), err)))
		return c.String(http.StatusBadRequest, fmt.Sprintf("Error: Unknown error, unable to read request %s", err))
	}

	var shortURLs []string
	err = json.Unmarshal(body, &shortURLs)
	if err != nil {
		h.logger.Error("StatusBadRequest: unknown error", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), err)))
		return c.String(http.StatusBadRequest, "Error: Unknown error, unable to read request")
	}

	userID, ok := c.Get("userID").(string)
	if !ok {
		h.logger.Error("Internal server error", zap.Error(urlErr.ErrUnableToGetUserIDFromContext))
		return c.NoContent(http.StatusInternalServerError)
	}

	workerPool := workerpool.NewWorkerPool(100, h.logger)
	workerPool.Start()
	defer workerPool.Stop()

	h.wg.Add(len(shortURLs))
	for _, shortURL := range shortURLs {
		log.Info("Submitting task", zap.String("delete shortURL", shortURL))
		url := shortURL
		userID := userID
		workerPool.Submit(func() error {
			defer h.wg.Done()
			return h.urlService.DeleteURLByUserID(context.WithoutCancel(c.Request().Context()), userID, url)
		})
	}

	return c.NoContent(http.StatusAccepted)
}

// AddBatch handles the addition of a batch of URLs.
func (h *URLHandler) AddBatch(c echo.Context) error {
	header := c.Request().Header.Get("Content-Type")
	if header != "application/json" {
		msg := "Content-Type header is not application/json"
		h.logger.Error("StatusUnsupportedMediaType: " + msg)
		return c.String(http.StatusUnsupportedMediaType, msg)
	}

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		h.logger.Error("StatusBadRequest: unknown error", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), err)))
		return c.String(http.StatusBadRequest, fmt.Sprintf("Error: Unknown error, unable to read request %s", err))
	}

	var urlBatchRequest []dto.URLBatchRequest
	err = json.Unmarshal(body, &urlBatchRequest)
	if err != nil {
		h.logger.Error("StatusBadRequest: unknown error", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), err)))
		return c.String(http.StatusBadRequest, "Error: Unknown error, unable to read request")
	}

	if len(urlBatchRequest) == 0 {
		h.logger.Error("StatusBadRequest: empty batch request", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), err)))
		return c.String(http.StatusBadRequest, "Error: empty batch request")
	}

	userID := c.Get("userID").(string)
	savedURLs, err := h.urlService.AddAll(c.Request().Context(), urlBatchRequest, h.urlPrefix, userID)
	if err != nil {
		h.logger.Error("StatusInternalServerError: unknown error", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), err)))
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unknown error: %s", err))
	}

	return c.JSON(http.StatusCreated, savedURLs)
}

// AddShorten handles the addition of a single URL (got as json).
func (h *URLHandler) AddShorten(c echo.Context) error {
	header := c.Request().Header.Get("Content-Type")
	if header != "application/json" {
		msg := "Content-Type header is not application/json"
		h.logger.Error("StatusUnsupportedMediaType: " + msg)
		return c.String(http.StatusUnsupportedMediaType, msg)
	}

	body, readBodyErr := io.ReadAll(c.Request().Body)
	if readBodyErr != nil {
		h.logger.Error("StatusBadRequest: unknown error", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), readBodyErr)))
		return c.String(http.StatusBadRequest, fmt.Sprintf("Error: Unknown error, unable to read request %s", readBodyErr))
	}

	var urlRequest dto.URLRequest
	unmarshalErr := json.Unmarshal(body, &urlRequest)
	if unmarshalErr != nil {
		h.logger.Error("StatusBadRequest: unknown error", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), unmarshalErr)))
		return c.String(http.StatusBadRequest, "Error: Unknown error, unable to read request")
	}

	if err := h.checkRequest(urlRequest.URL); err != nil {
		h.logger.Error("StatusBadRequest: unable to handle empty request", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), err)))
		return c.String(http.StatusBadRequest, "Error: Unable to handle empty request")
	}

	userID, ok := c.Get("userID").(string)
	if !ok {
		h.logger.Error("Internal server error", zap.Error(urlErr.ErrUnableToGetUserIDFromContext))
		return c.NoContent(http.StatusInternalServerError)
	}

	url, err := h.urlService.Add(c.Request().Context(), urlRequest.URL, h.urlPrefix, userID)
	if err != nil && !errors.Is(err, urlErr.ErrURLAlreadyExists) {
		h.logger.Error("StatusBadRequest: unknown error", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), err)))
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unknown error: %s", err))
	}

	response := &dto.URLResponse{
		Result: url.Shortened,
	}

	if errors.Is(err, urlErr.ErrURLAlreadyExists) {
		h.logger.Warn("StatusConflict: url already exists", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), err)))
		return c.JSON(http.StatusConflict, response)
	}

	return c.JSON(http.StatusCreated, response)
}

// AddURL handles the addition of a URL (got as plain text).
func (h *URLHandler) AddURL(c echo.Context) error {
	body, readErr := io.ReadAll(c.Request().Body)
	if readErr != nil {
		h.logger.Error("StatusBadRequest: unknown error", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), readErr)))
		return c.String(http.StatusBadRequest, "Error: Unknown error, unable to read request")
	}

	if err := h.checkRequest(string(body)); err != nil {
		h.logger.Error("StatusBadRequest: unable to handle empty request", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), err)))
		return c.String(http.StatusBadRequest, "Error: Unable to handle empty request")
	}

	userID, ok := c.Get("userID").(string)
	if !ok {
		h.logger.Error("Internal server error", zap.Error(urlErr.ErrUnableToGetUserIDFromContext))
		return c.NoContent(http.StatusInternalServerError)
	}

	url, err := h.urlService.Add(c.Request().Context(), string(body), h.urlPrefix, userID)
	if err != nil && !errors.Is(err, urlErr.ErrURLAlreadyExists) {
		h.logger.Error("StatusInternalServerError: Unknown error:", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), err)))
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unknown error: %s", err))
	}

	if errors.Is(err, urlErr.ErrURLAlreadyExists) {
		h.logger.Warn("StatusConflict: url already exists", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), err)))
		return c.String(http.StatusConflict, url.Shortened)
	}

	c.Response().WriteHeader(http.StatusCreated)
	return c.String(http.StatusCreated, url.Shortened)
}

// ClearAll deletes all data and returns an error if any.
//
// Deletes all saved urls.
func (h *URLHandler) ClearAll(c echo.Context) error {
	if err := h.urlService.DeleteAll(c.Request().Context()); err != nil {
		h.logger.Error("StatusInternalServerError: Unknown error:", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), err)))
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unknown error: %s", err))
	}

	return c.String(http.StatusOK, "All data deleted")
}

// FindAll retrieves all URLs.
//
// Retrieves all saved urls.
func (h *URLHandler) FindAll(c echo.Context) error {
	urls, err := h.urlService.GetAll(c.Request().Context())
	if err != nil {
		h.logger.Error("StatusInternalServerError: Unknown error:", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), err)))
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unknown error: %s", err))
	}

	return c.String(http.StatusOK, strings.Join(urls, ", "))
}

// FindURL finds the URL based on the given ID from echo context.
//
// Finds the URL based on the given ID.
func (h *URLHandler) FindURL(c echo.Context) error {
	id := (strings.Split(c.Request().URL.Path, "/"))[1]

	if err := h.checkRequest(id); err != nil {
		h.logger.Error("StatusBadRequest: Unable to handle empty request", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), err)))
		return c.String(http.StatusBadRequest, "Error: Unable to handle empty request")
	}

	originalURL, err := h.urlService.GetByyID(c.Request().Context(), id)

	var message string
	var status int

	switch {
	case errors.Is(err, urlErr.ErrURLNotFound):
		h.logger.Info("StatusBadRequest: url not found", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), err)))
		status = http.StatusBadRequest
		message = fmt.Sprintf("URL with id %s not found", id)

	case errors.Is(err, urlErr.ErrURLDeleted):
		h.logger.Info("StatusBadRequest: url not found", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), err)))
		status = http.StatusGone
		message = fmt.Sprintf("URL with id %s has been deleted", id)

	case err != nil:
		h.logger.Error("InternalServerError", zap.Error(fmt.Errorf("%s %w", apperr.Caller(), err)))
		status = http.StatusInternalServerError
		message = fmt.Sprintf("Unknown error: %s", err)

	default:
		c.Response().Header().Set("Location", originalURL)
		status = http.StatusTemporaryRedirect
		message = ""
	}

	return c.String(status, message)
}

// checkRequest checks if the request is empty.
func (h *URLHandler) checkRequest(s string) error {
	if len(s) == 0 {
		return apperr.NewValueError("Unable to handle empty request", apperr.Caller(), urlErr.ErrEmptyRequest)
	}

	return nil
}

// Ping is a function that handles the ping request.
//
// Check database connection.
func (h *URLHandler) Ping(c echo.Context) error {
	status := http.StatusOK
	err := h.urlService.Ping(c.Request().Context())
	if err != nil {
		status = http.StatusInternalServerError
	}

	return c.NoContent(status)
}
