package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"

	"github.com/msmkdenis/yap-shortener/internal/apperrors"
	"github.com/msmkdenis/yap-shortener/internal/handlers/dto"
	"github.com/msmkdenis/yap-shortener/internal/middleware"
	"github.com/msmkdenis/yap-shortener/internal/model"
	"github.com/msmkdenis/yap-shortener/internal/utils"
)

type URLHandler struct {
	urlService URLService
	urlPrefix  string
	jwtManager *utils.JWTManager
	logger     *zap.Logger
}

type URLService interface {
	Add(ctx context.Context, s string, host string, userID string) (*model.URL, error)
	AddAll(ctx context.Context, urls []dto.URLBatchRequest, host string, userID string) ([]dto.URLBatchResponse, error)
	GetAll(ctx context.Context) ([]string, error)
	GetAllByUserID(ctx context.Context, userID string) ([]dto.URLBatchResponseByUserID, error)
	DeleteAll(ctx context.Context) error
	DeleteAllByUserID(ctx context.Context, userID string, shortURLs []string) error
	GetByyID(ctx context.Context, key string) (string, error)
	Ping(ctx context.Context) error
}

func NewURLHandler(e *echo.Echo, service URLService, urlPrefix string, jwtManager *utils.JWTManager, logger *zap.Logger) *URLHandler {
	handler := &URLHandler{
		urlService: service,
		urlPrefix:  urlPrefix,
		jwtManager: jwtManager,
		logger:     logger,
	}

	requestLogger := middleware.InitRequestLogger(logger)
	jwtCheckerCreator := middleware.InitJWTCheckerCreator(jwtManager, logger)
	jwtAuth := middleware.InitJWTAuth(jwtManager, logger)

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

	return handler
}

func (h *URLHandler) FindAllURLByUserID(c echo.Context) error {
	userID := c.Get("userID").(string)
	savedURLs, err := h.urlService.GetAllByUserID(c.Request().Context(), userID)
	if err != nil && !errors.Is(err, apperrors.ErrURLNotFound) {
		h.logger.Error("StatusBadRequest: unknown error", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), err)))
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unknown error: %s", err))
	}

	if errors.Is(err, apperrors.ErrURLNotFound) {
		h.logger.Warn("StatusNoContent: urls not found", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), err)))
		return c.NoContent(http.StatusNoContent)
	}

	return c.JSON(http.StatusOK, savedURLs)
}

func (h *URLHandler) DeleteAllURLsByUserID(c echo.Context) error {
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

	var shortURLs []string
	err = json.Unmarshal(body, &shortURLs)
	if err != nil {
		h.logger.Error("StatusBadRequest: unknown error", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), err)))
		return c.String(http.StatusBadRequest, "Error: Unknown error, unable to read request")
	}

	workerPool := utils.NewWorkerPool(100, h.logger)
	workerPool.Start()

	for _, shortURL := range shortURLs {
		log.Info("Submitting task", zap.String("delete shortURL", shortURL))
		url := []string{shortURL}
		workerPool.Submit(func() {
			err = h.urlService.DeleteAllByUserID(c.Request().Context(), c.Get("userID").(string), url)
			if err != nil && !errors.Is(err, apperrors.ErrURLNotFound) {
				h.logger.Error("StatusBadRequest: unknown error", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), err)))
			}
		})
	}

	defer workerPool.Stop()
	workerPool.Wait() // also we can skip this line and return 202 since there is no need in returning info to client

	return c.NoContent(http.StatusAccepted)
}

func (h *URLHandler) AddBatch(c echo.Context) error {
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

	var urlBatchRequest []dto.URLBatchRequest
	err = json.Unmarshal(body, &urlBatchRequest)
	if err != nil {
		h.logger.Error("StatusBadRequest: unknown error", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), err)))
		return c.String(http.StatusBadRequest, "Error: Unknown error, unable to read request")
	}

	if len(urlBatchRequest) == 0 {
		h.logger.Error("StatusBadRequest: empty batch request", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), err)))
		return c.String(http.StatusBadRequest, "Error: empty batch request")
	}

	userID := c.Get("userID").(string)
	savedURLs, err := h.urlService.AddAll(c.Request().Context(), urlBatchRequest, h.urlPrefix, userID)
	if err != nil {
		h.logger.Error("StatusInternalServerError: unknown error", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), err)))
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unknown error: %s", err))
	}

	return c.JSON(http.StatusCreated, savedURLs)
}

func (h *URLHandler) AddShorten(c echo.Context) error {
	header := c.Request().Header.Get("Content-Type")
	if header != "application/json" {
		msg := "Content-Type header is not application/json"
		h.logger.Error("StatusUnsupportedMediaType: " + msg)
		return c.String(http.StatusUnsupportedMediaType, msg)
	}

	body, readBodyErr := io.ReadAll(c.Request().Body)
	if readBodyErr != nil {
		h.logger.Error("StatusBadRequest: unknown error", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), readBodyErr)))
		return c.String(http.StatusBadRequest, fmt.Sprintf("Error: Unknown error, unable to read request %s", readBodyErr))
	}

	var urlRequest dto.URLRequest
	unmarshalErr := json.Unmarshal(body, &urlRequest)
	if unmarshalErr != nil {
		h.logger.Error("StatusBadRequest: unknown error", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), unmarshalErr)))
		return c.String(http.StatusBadRequest, "Error: Unknown error, unable to read request")
	}

	if err := h.checkRequest(urlRequest.URL); err != nil {
		h.logger.Error("StatusBadRequest: unable to handle empty request", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), err)))
		return c.String(http.StatusBadRequest, "Error: Unable to handle empty request")
	}

	userID := c.Get("userID").(string)
	url, err := h.urlService.Add(c.Request().Context(), urlRequest.URL, h.urlPrefix, userID)
	if err != nil && !errors.Is(err, apperrors.ErrURLAlreadyExists) {
		h.logger.Error("StatusBadRequest: unknown error", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), err)))
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unknown error: %s", err))
	}

	response := &dto.URLResponse{
		Result: url.Shortened,
	}

	if errors.Is(err, apperrors.ErrURLAlreadyExists) {
		h.logger.Warn("StatusConflict: url already exists", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), err)))
		return c.JSON(http.StatusConflict, response)
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *URLHandler) AddURL(c echo.Context) error {
	body, readErr := io.ReadAll(c.Request().Body)
	if readErr != nil {
		h.logger.Error("StatusBadRequest: unknown error", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), readErr)))
		return c.String(http.StatusBadRequest, "Error: Unknown error, unable to read request")
	}

	if err := h.checkRequest(string(body)); err != nil {
		h.logger.Error("StatusBadRequest: unable to handle empty request", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), err)))
		return c.String(http.StatusBadRequest, "Error: Unable to handle empty request")
	}

	userID := c.Get("userID").(string)
	url, err := h.urlService.Add(c.Request().Context(), string(body), h.urlPrefix, userID)
	if err != nil && !errors.Is(err, apperrors.ErrURLAlreadyExists) {
		h.logger.Error("StatusInternalServerError: Unknown error:", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), err)))
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unknown error: %s", err))
	}

	if errors.Is(err, apperrors.ErrURLAlreadyExists) {
		h.logger.Warn("StatusConflict: url already exists", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), err)))
		return c.String(http.StatusConflict, url.Shortened)
	}

	c.Response().WriteHeader(http.StatusCreated)
	return c.String(http.StatusCreated, url.Shortened)
}

func (h *URLHandler) ClearAll(c echo.Context) error {
	if err := h.urlService.DeleteAll(c.Request().Context()); err != nil {
		h.logger.Error("StatusInternalServerError: Unknown error:", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), err)))
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unknown error: %s", err))
	}

	return c.String(http.StatusOK, "All data deleted")
}

func (h *URLHandler) FindAll(c echo.Context) error {
	urls, err := h.urlService.GetAll(c.Request().Context())
	if err != nil {
		h.logger.Error("StatusInternalServerError: Unknown error:", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), err)))
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unknown error: %s", err))
	}

	return c.String(http.StatusOK, strings.Join(urls, ", "))
}

func (h *URLHandler) FindURL(c echo.Context) error {
	id := (strings.Split(c.Request().URL.Path, "/"))[1]

	if err := h.checkRequest(id); err != nil {
		h.logger.Error("StatusBadRequest: Unable to handle empty request", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), err)))
		return c.String(http.StatusBadRequest, "Error: Unable to handle empty request")
	}

	originalURL, err := h.urlService.GetByyID(c.Request().Context(), id)

	var message string
	var status int

	switch {
	case errors.Is(err, apperrors.ErrURLNotFound):
		h.logger.Info("StatusBadRequest: url not found", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), err)))
		status = http.StatusBadRequest
		message = fmt.Sprintf("URL with id %s not found", id)

	case errors.Is(err, apperrors.ErrURLDeleted):
		h.logger.Info("StatusBadRequest: url not found", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), err)))
		status = http.StatusGone
		message = fmt.Sprintf("URL with id %s has been deleted", id)

	case err != nil:
		h.logger.Error("InternalServerError", zap.Error(fmt.Errorf("caller: %s %w", utils.Caller(), err)))
		status = http.StatusInternalServerError
		message = fmt.Sprintf("Unknown error: %s", err)

	default:
		c.Response().Header().Set("Location", originalURL)
		status = http.StatusTemporaryRedirect
		message = ""
	}

	return c.String(status, message)
}

func (h *URLHandler) checkRequest(s string) error {
	if len(s) == 0 {
		return apperrors.NewValueError("Unable to handle empty request", utils.Caller(), apperrors.ErrEmptyRequest)
	}

	return nil
}

func (h *URLHandler) Ping(c echo.Context) error {
	status := http.StatusOK
	err := h.urlService.Ping(c.Request().Context())
	if err != nil {
		status = http.StatusInternalServerError
	}

	return c.NoContent(status)
}
