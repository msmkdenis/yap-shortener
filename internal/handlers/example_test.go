package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/msmkdenis/yap-shortener/internal/config"
	"github.com/msmkdenis/yap-shortener/internal/dto"
	"github.com/msmkdenis/yap-shortener/internal/repository/memory"
	"github.com/msmkdenis/yap-shortener/internal/service"
	"github.com/msmkdenis/yap-shortener/pkg/echopprof"
	"github.com/msmkdenis/yap-shortener/pkg/jwtgen"
)

var cfgExampleTest = &config.Config{
	URLServer:       "8080",
	URLPrefix:       "http://localhost:8080",
	FileStoragePath: "/tmp/short-url-db-test.json",
	TokenName:       "test",
	SecretKey:       "test",
}

func ExampleURLHandler_AddURL() {
	cfg := cfgExampleTest
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Unable to initialize zap logger", zap.Error(err))
	}
	jwtManager := jwtgen.InitJWTManager(cfg.TokenName, cfg.SecretKey, logger)
	repository := memory.NewURLRepository(logger)
	urlService := service.NewURLService(repository, logger)

	e := echo.New()
	echopprof.Wrap(e)
	h := NewURLHandler(e, urlService, cfg.URLPrefix, jwtManager, logger)

	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/", strings.NewReader("https://example.com"))
	w := httptest.NewRecorder()
	l := e.NewContext(request, w)
	l.Set("userID", "token")

	err = h.AddURL(l)

	answer := w.Body.String()
	fmt.Println(answer)
	fmt.Println(w.Code)
	// Output: http://localhost:8080/Yzk4NGQ
	// 201
}

func ExampleURLHandler_FindAllURLByUserID() {
	cfg := cfgExampleTest
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Unable to initialize zap logger", zap.Error(err))
	}
	jwtManager := jwtgen.InitJWTManager(cfg.TokenName, cfg.SecretKey, logger)
	repository := memory.NewURLRepository(logger)
	urlService := service.NewURLService(repository, logger)

	e := echo.New()
	echopprof.Wrap(e)
	h := NewURLHandler(e, urlService, cfg.URLPrefix, jwtManager, logger)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/user/urls", strings.NewReader(""))
	w := httptest.NewRecorder()
	l := e.NewContext(request, w)
	request.Header.Set("Content-Type", "application/json")
	l.Set("userID", "token")

	urlService.Add(context.Background(), "https://example.com", "localhost:8080", "token")
	urlService.Add(context.Background(), "https://new.com", "localhost:8080", "token")

	h.FindAllURLByUserID(l)

	fmt.Println(w.Code)
	// Output: 200
}

func ExampleURLHandler_AddBatch() {
	cfg := cfgExampleTest
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Unable to initialize zap logger", zap.Error(err))
	}
	jwtManager := jwtgen.InitJWTManager(cfg.TokenName, cfg.SecretKey, logger)
	repository := memory.NewURLRepository(logger)
	urlService := service.NewURLService(repository, logger)

	e := echo.New()
	echopprof.Wrap(e)
	h := NewURLHandler(e, urlService, cfg.URLPrefix, jwtManager, logger)

	fullURL1 := dto.URLBatchRequest{
		CorrelationID: "2d5d144a-f272-40d3-b3aa-d4b1b9da277c",
		OriginalURL:   "https://example.com",
	}
	fullURL2 := dto.URLBatchRequest{
		CorrelationID: "8a0aea8b-c62e-4f2f-a4b6-c97b9bdf2e99",
		OriginalURL:   "https://new.com",
	}
	requestBody := []dto.URLBatchRequest{
		fullURL1,
		fullURL2,
	}

	body, _ := json.Marshal(requestBody)

	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/shorten/batch", strings.NewReader(string(body)))
	w := httptest.NewRecorder()
	l := e.NewContext(request, w)
	request.Header.Set("Content-Type", "application/json")
	l.Set("userID", "token")

	h.AddBatch(l)

	answer := w.Body.String()
	fmt.Println(answer)
	fmt.Println(w.Code)
	// Output: [{"correlation_id":"2d5d144a-f272-40d3-b3aa-d4b1b9da277c","short_url":"http://localhost:8080/Yzk4NGQ"},{"correlation_id":"8a0aea8b-c62e-4f2f-a4b6-c97b9bdf2e99","short_url":"http://localhost:8080/YTFjMTY"}]
	//
	// 201
}

func ExampleURLHandler_FindAll() {
	cfg := cfgExampleTest
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Unable to initialize zap logger", zap.Error(err))
	}
	jwtManager := jwtgen.InitJWTManager(cfg.TokenName, cfg.SecretKey, logger)
	repository := memory.NewURLRepository(logger)
	urlService := service.NewURLService(repository, logger)

	e := echo.New()
	echopprof.Wrap(e)
	h := NewURLHandler(e, urlService, cfg.URLPrefix, jwtManager, logger)

	fullURL1 := dto.URLBatchRequest{
		CorrelationID: "2d5d144a-f272-40d3-b3aa-d4b1b9da277c",
		OriginalURL:   "https://example.com",
	}
	fullURL2 := dto.URLBatchRequest{
		CorrelationID: "8a0aea8b-c62e-4f2f-a4b6-c97b9bdf2e99",
		OriginalURL:   "https://new.com",
	}
	requestBody := []dto.URLBatchRequest{
		fullURL1,
		fullURL2,
	}

	urlService.AddAll(context.Background(), requestBody, "localhost:8080", "token")

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", strings.NewReader(""))
	w := httptest.NewRecorder()
	l := e.NewContext(request, w)
	l.Set("userID", "token")

	h.FindAll(l)

	fmt.Println(w.Code)
	// Output: 200
}

func ExampleURLHandler_FindURL() {
	cfg := cfgExampleTest
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Unable to initialize zap logger", zap.Error(err))
	}
	jwtManager := jwtgen.InitJWTManager(cfg.TokenName, cfg.SecretKey, logger)
	repository := memory.NewURLRepository(logger)
	urlService := service.NewURLService(repository, logger)

	e := echo.New()
	echopprof.Wrap(e)
	h := NewURLHandler(e, urlService, cfg.URLPrefix, jwtManager, logger)

	urlService.Add(context.Background(), "https://example.com", "localhost:8080", "token")

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/Yzk4NGQ", strings.NewReader(""))
	w := httptest.NewRecorder()
	l := e.NewContext(request, w)
	l.Set("userID", "token")

	h.FindURL(l)

	fmt.Println(w.Code)
	// Output: 307
}

func ExampleURLHandler_AddShorten() {
	cfg := cfgExampleTest
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Unable to initialize zap logger", zap.Error(err))
	}
	jwtManager := jwtgen.InitJWTManager(cfg.TokenName, cfg.SecretKey, logger)
	repository := memory.NewURLRepository(logger)
	urlService := service.NewURLService(repository, logger)

	e := echo.New()
	echopprof.Wrap(e)
	h := NewURLHandler(e, urlService, cfg.URLPrefix, jwtManager, logger)

	requestBody := dto.URLRequest{URL: "https://example.com"}
	body, _ := json.Marshal(requestBody)

	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/shorten", strings.NewReader(string(body)))
	w := httptest.NewRecorder()
	l := e.NewContext(request, w)
	request.Header.Set("Content-Type", "application/json")
	l.Set("userID", "token")

	h.AddShorten(l)

	answer := w.Body.String()
	fmt.Println(answer)
	fmt.Println(w.Code)
	// Output: {"result":"http://localhost:8080/Yzk4NGQ"}
	//
	// 201
}
