package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/msmkdenis/yap-shortener/internal/config"
	"github.com/msmkdenis/yap-shortener/internal/repository/memory"
	"github.com/msmkdenis/yap-shortener/internal/service"
	"github.com/msmkdenis/yap-shortener/internal/utils"
	"github.com/msmkdenis/yap-shortener/pkg/echopprof"
)

func ExampleURLHandler_AddURL() {
	cfg := *config.NewConfig()
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Unable to initialize zap logger", zap.Error(err))
	}
	jwtManager := utils.InitJWTManager(cfg.TokenName, cfg.SecretKey, logger)
	repository := memory.NewURLRepository(logger)
	urlService := service.NewURLService(repository, logger)

	e := echo.New()
	echopprof.Wrap(e)
	h := NewURLHandler(e, urlService, cfg.URLPrefix, jwtManager, logger)

	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/", strings.NewReader("https://example.com"))
	w := httptest.NewRecorder()
	l := e.NewContext(request, w)
	request.Header.Set("Content-Type", "application/json")
	l.Set("userID", "token")
	l.Set("userID", "token")

	h.AddURL(l)

	answer := w.Body.String()
	fmt.Println(answer)
	// Output: http://localhost:8080/Yzk4NGQ
}
