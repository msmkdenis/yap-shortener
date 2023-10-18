package handlers

import (
	"github.com/labstack/echo/v4"
	storage2 "github.com/msmkdenis/yap-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestURLHandler(t *testing.T) {

	storage2.GlobalRepository = storage2.NewMemoryRepository()

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

	e := echo.New()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			switch test.method {

			case http.MethodGet:
				preRequest := httptest.NewRequest(http.MethodPost, "http://localhost:8080/", strings.NewReader(test.body))
				preW := httptest.NewRecorder()
				c := e.NewContext(preRequest, preW)
				PostURL(c)

				request := httptest.NewRequest(test.method, test.path, strings.NewReader(test.body))
				w := httptest.NewRecorder()
				b := e.NewContext(request, w)
				GetURL(b)
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
				PostURL(l)
				res := w.Result()
				assert.Equal(t, test.want.code, res.StatusCode)
				defer res.Body.Close()
				_, err := io.ReadAll(res.Body)
				require.NoError(t, err)

			}
		})
	}
}