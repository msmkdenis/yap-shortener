package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/msmkdenis/yap-shortener/internal/config"
	"github.com/msmkdenis/yap-shortener/internal/repository/memory"
	"github.com/msmkdenis/yap-shortener/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var cfgMock = &config.Config{
	URLServer: "8080",
	URLPrefix: "http://localhost:8080",
}

func TestURLHandler(t *testing.T) {
	//cfg := *config.NewConfig()
	urlRepository := memory.NewURLRepository()
	urlService := service.NewURLService(urlRepository)
	logger, _ := zap.NewProduction()

	e := echo.New()

	h := New(e, urlService, cfgMock.URLPrefix, logger)

	type want struct {
		code     int
		location string
	}

	tests := []struct {
		name   string
		method string
		body   string
		path   string
		want   want
	}{
		{
			name:   "positive POST test #1",
			method: http.MethodPost,
			body:   "https://practicum.yandex.ru/",
			path:   "http://localhost:8080/",
			want: want{
				code: http.StatusCreated,
			},
		},
		{
			name:   "positive POST test #2",
			method: http.MethodPost,
			body:   "https://ru.tradingview.com/",
			path:   "http://localhost:8080/",
			want: want{
				code: http.StatusCreated,
			},
		},
		{
			name:   "negative POST test #2 (empty request)",
			method: http.MethodPost,
			body:   "",
			path:   "http://localhost:8080/",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "positive GET test #1",
			method: http.MethodGet,
			body:   "https://practicum.yandex.ru/",
			path:   "http://localhost:8080/MGRkMTk",
			want: want{
				code:     http.StatusTemporaryRedirect,
				location: "https://practicum.yandex.ru",
			},
		},
		{
			name:   "negative GET test #1 (wrong id)",
			method: http.MethodGet,
			body:   "https://practicum.yandex.ru/",
			path:   "http://localhost:8080/MGRkiuY",
			want: want{
				code:     http.StatusBadRequest,
				location: "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			switch test.method {

			case http.MethodGet:
				preRequest := httptest.NewRequest(http.MethodPost, "http://localhost:8080/", strings.NewReader(test.body))
				preW := httptest.NewRecorder()
				c := e.NewContext(preRequest, preW)
				h.PostURL(c)

				request := httptest.NewRequest(test.method, test.path, strings.NewReader(test.body))
				w := httptest.NewRecorder()
				b := e.NewContext(request, w)
				h.GetURL(b)
				res := w.Result()
				defer res.Body.Close()
				assert.Equal(t, test.want.code, res.StatusCode)
				if res.StatusCode == http.StatusBadRequest {
					assert.Equal(t, res.Header.Get("Location"), "")
				}

			case http.MethodPost:
				request := httptest.NewRequest(test.method, test.path, strings.NewReader(test.body))
				w := httptest.NewRecorder()
				l := e.NewContext(request, w)
				h.PostURL(l)
				res := w.Result()
				assert.Equal(t, test.want.code, res.StatusCode)
				defer res.Body.Close()
				_, err := io.ReadAll(res.Body)
				require.NoError(t, err)

			}
		})
	}
}

func TestPostShorten(t *testing.T) {
	//cfg := *config.NewConfig()
	urlRepository := memory.NewURLRepository()
	urlService := service.NewURLService(urlRepository)
	logger, _ := zap.NewProduction()

	e := echo.New()

	h := New(e, urlService, cfgMock.URLPrefix, logger)

	type want struct {
		code     int
		response string
	}

	tests := []struct {
		name        string
		method      string
		body        string
		contentType string
		path        string
		want        want
	}{
		{
			name:        "positive PostShorten test #1",
			method:      http.MethodPost,
			body:        `{"url":"https://practicum.yandex.ru/"}`,
			contentType: "application/json",
			path:        "http://localhost:8080/api/shorten",
			want: want{
				code:     http.StatusCreated,
				response: `{"result":"http://localhost:8080/MGRkMTk"}` + "\n",
			},
		},
		{
			name:        "negative PostShorten test #1",
			method:      http.MethodPost,
			body:        `{"url":"https://practicum.yandex.ru/"}`,
			contentType: "",
			path:        "http://localhost:8080/api/shorten",
			want: want{
				code:     http.StatusUnsupportedMediaType,
				response: "Content-Type header is not application/json",
			},
		},
		{
			name:        "negative PostShorten test #2",
			method:      http.MethodPost,
			body:        `{}`,
			contentType: "application/json",
			path:        "http://localhost:8080/api/shorten",
			want: want{
				code:     http.StatusBadRequest,
				response: "Error: Unable to handle empty body",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			switch test.method {
			case http.MethodPost:
				request := httptest.NewRequest(test.method, test.path, strings.NewReader(test.body))
				request.Header.Set("Content-Type", test.contentType)
				w := httptest.NewRecorder()
				l := e.NewContext(request, w)
				h.PostShorten(l)
				res := w.Result()
				assert.Equal(t, test.want.code, res.StatusCode)
				defer res.Body.Close()
				response, err := io.ReadAll(res.Body)
				assert.Equal(t, test.want.response, string(response))
				require.NoError(t, err)
			}
		})
	}
}
