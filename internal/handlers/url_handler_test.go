package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/msmkdenis/yap-shortener/internal/repository/memory"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"

	"github.com/msmkdenis/yap-shortener/internal/apperrors"
	"github.com/msmkdenis/yap-shortener/internal/config"
	mock "github.com/msmkdenis/yap-shortener/internal/mocks"
	"github.com/msmkdenis/yap-shortener/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var cfgMock = &config.Config{
	URLServer:       "8080",
	URLPrefix:       "http://localhost:8080",
	FileStoragePath: "/tmp/short-url-db-test.json",
}

func TestURLHandler(t *testing.T) {
	logger, _ := zap.NewProduction()
	urlRepository := memory.NewURLRepository(logger)
	urlService := service.NewURLService(urlRepository, logger)

	e := echo.New()

	h := NewURLHandler(e, urlService, cfgMock.URLPrefix, logger)

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
			body:   "http://fgdgdfg.com/qwpoeipqowei",
			path:   "http://localhost:8080/",
			want: want{
				code: http.StatusCreated,
			},
		},
		{
			name:   "positive POST test #2",
			method: http.MethodPost,
			body:   "http://ip95f7tnksykxx.biz/q3jl16cuadw/viajydc/kp8rl2",
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
			body:   "https://web.telegram.org/a/",
			path:   "http://localhost:8080/ZDFiODN",
			want: want{
				code:     http.StatusTemporaryRedirect,
				location: "https://web.telegram.org/a/",
			},
		},
		{
			name:   "negative GET test #1 (wrong id)",
			method: http.MethodGet,
			body:   "https://stackoverflow.com/",
			path:   "http://localhost:8080/918237918273",
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
				err := h.AddURL(c)
				require.NoError(t, err)

				request := httptest.NewRequest(test.method, test.path, strings.NewReader(test.body))
				w := httptest.NewRecorder()
				b := e.NewContext(request, w)
				err = h.FindURL(b)
				require.NoError(t, err)
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
				err := h.AddURL(l)
				require.NoError(t, err)
				res := w.Result()
				assert.Equal(t, test.want.code, res.StatusCode)
				defer res.Body.Close()
				_, err = io.ReadAll(res.Body)
				require.NoError(t, err)
			}
		})
	}
}

func TestPostShorten(t *testing.T) {
	logger, _ := zap.NewProduction()
	urlRepository := memory.NewURLRepository(logger)
	urlService := service.NewURLService(urlRepository, logger)

	e := echo.New()

	h := NewURLHandler(e, urlService, cfgMock.URLPrefix, logger)

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
			body:        `{"url":"https://www.dns-shop.ru/"}`,
			contentType: "application/json",
			path:        "http://localhost:8080/api/shorten",
			want: want{
				code:     http.StatusCreated,
				response: `{"result":"http://localhost:8080/ZmM0NTQ"}` + "\n",
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
				response: "Error: Unable to handle empty request",
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
				_ = h.AddShorten(l)
				res := w.Result()
				assert.Equal(t, test.want.code, res.StatusCode)
				defer res.Body.Close()
				response, err := io.ReadAll(res.Body)
				assert.Equal(t, test.want.response, string(response))
				require.NoError(t, err)
			}
		})
	}
	_ = urlService.DeleteAll
}

func TestGetURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock.NewMockURLService(ctrl)

	message := "https://practicum.yandex.ru/"

	s.EXPECT().GetByyID(gomock.Any(), gomock.Any()).AnyTimes().Return(message, nil)

	logger, _ := zap.NewProduction()
	e := echo.New()
	h := NewURLHandler(e, s, cfgMock.URLPrefix, logger)
	defer e.Close()

	testCases := []struct {
		name             string
		method           string
		body             string
		expectedCode     int
		path             string
		expectedBody     string
		expectedLocation string
	}{
		{
			name:             "BadRequest - ID is empty",
			method:           http.MethodGet,
			expectedCode:     http.StatusBadRequest,
			path:             "http://localhost:8080/",
			expectedBody:     "Error: Unable to handle empty request",
			expectedLocation: "",
		},
		{
			name:             "TemporaryRedirect - ID is not empty",
			method:           http.MethodGet,
			expectedCode:     http.StatusTemporaryRedirect,
			path:             "http://localhost:8080/MGRkMTk",
			expectedBody:     "",
			expectedLocation: "https://practicum.yandex.ru/",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.path, strings.NewReader(""))
			w := httptest.NewRecorder()
			l := e.NewContext(request, w)
			err := h.FindURL(l)
			require.NoError(t, err)
			res := w.Result()
			assert.Equal(t, test.expectedCode, res.StatusCode)
			assert.Equal(t, test.expectedLocation, res.Header.Get("Location"))
			response, err := io.ReadAll(res.Body)
			res.Body.Close()
			assert.Equal(t, test.expectedBody, string(response))
			require.NoError(t, err)
		})
	}

}

func TestGetURLError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock.NewMockURLService(ctrl)

	logger, _ := zap.NewProduction()
	e := echo.New()
	h := NewURLHandler(e, s, cfgMock.URLPrefix, logger)
	defer e.Close()

	s.EXPECT().GetByyID(gomock.Any(), gomock.Any()).Times(1).Return("", apperrors.ErrURLNotFound)

	testCaseWithError := []struct {
		name         string
		method       string
		body         string
		expectedCode int
		path         string
		expectedBody string
	}{
		{
			name:         "BadRequest - url not found",
			method:       http.MethodGet,
			expectedCode: http.StatusBadRequest,
			path:         "http://localhost:8080/MGRkMTk",
			expectedBody: "URL with id MGRkMTk not found",
		},
	}

	for _, test := range testCaseWithError {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.path, strings.NewReader(""))
			w := httptest.NewRecorder()
			l := e.NewContext(request, w)
			err := h.FindURL(l)
			require.NoError(t, err)
			res := w.Result()
			assert.Equal(t, test.expectedCode, res.StatusCode)
			response, err := io.ReadAll(res.Body)
			res.Body.Close()
			assert.Equal(t, test.expectedBody, string(response))
			require.NoError(t, err)
		})
	}
}
