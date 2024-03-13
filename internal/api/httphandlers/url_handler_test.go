package httphandlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/msmkdenis/yap-shortener/internal/config"
	"github.com/msmkdenis/yap-shortener/internal/dto"
	"github.com/msmkdenis/yap-shortener/internal/middleware"
	mock "github.com/msmkdenis/yap-shortener/internal/mocks"
	"github.com/msmkdenis/yap-shortener/internal/model"
	urlErr "github.com/msmkdenis/yap-shortener/internal/urlerr"
	"github.com/msmkdenis/yap-shortener/pkg/jwtgen"
)

var cfgMock = &config.Config{
	URLServer:       "8080",
	URLPrefix:       "http://localhost:8080",
	TrustedSubnet:   "",
	FileStoragePath: "/tmp/short-url-db-test.json",
	TokenName:       "test",
	SecretKey:       "test",
}

type URLHandlerTestSuite struct {
	suite.Suite
	h          *URLShorten
	urlService *mock.MockURLService
	echo       *echo.Echo
	ctrl       *gomock.Controller
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(URLHandlerTestSuite))
}

func (s *URLHandlerTestSuite) SetupTest() {
	cfgMock.URLPrefix = "http://localhost:8080"
	logger, _ := zap.NewProduction()
	jwtManager := jwtgen.InitJWTManager(cfgMock.TokenName, cfgMock.SecretKey, logger)
	jwtCheckerCreator := middleware.InitJWTCheckerCreator(jwtManager, logger)
	jwtAuth := middleware.InitJWTAuth(jwtManager, logger)
	s.ctrl = gomock.NewController(s.T())
	s.echo = echo.New()
	s.urlService = mock.NewMockURLService(s.ctrl)
	s.h = NewURLShorten(s.echo, s.urlService, cfgMock.URLPrefix, cfgMock.TrustedSubnet, jwtCheckerCreator, jwtAuth, logger, &sync.WaitGroup{})
}

func (s *URLHandlerTestSuite) TestDeleteAllURLsByUserID_Unauthorized() {
	defer func(echo *echo.Echo) {
		err := echo.Close()
		assert.NoError(s.T(), err)
	}(s.echo)

	testCases := []struct {
		name         string
		method       string
		body         string
		expectedCode int
		path         string
		expectedBody string
	}{
		{
			name:         "BadRequest - unauthorized",
			method:       http.MethodDelete,
			expectedCode: http.StatusUnauthorized,
			path:         "http://localhost:8080/api/user/urls",
			expectedBody: "",
		},
	}

	for _, test := range testCases {
		s.T().Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.path, strings.NewReader(""))
			w := httptest.NewRecorder()
			s.echo.ServeHTTP(w, request)
			assert.Equal(t, test.expectedCode, w.Code)
		})
	}
}

func (s *URLHandlerTestSuite) TestDeleteAllURLsByUserID_WrongMediaType() {
	testCases := []struct {
		name         string
		method       string
		body         string
		expectedCode int
		path         string
		expectedBody string
	}{
		{
			name:         "BadRequest - unsupported media type",
			method:       http.MethodDelete,
			expectedCode: http.StatusUnsupportedMediaType,
			path:         "http://localhost:8080/api/user/urls",
			expectedBody: "",
		},
	}

	for _, test := range testCases {
		s.T().Run(test.name, func(t *testing.T) {
			s.urlService.EXPECT().DeleteURLByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			request := httptest.NewRequest(test.method, test.path, strings.NewReader(""))
			request.Header.Set("Content-Type", "")
			w := httptest.NewRecorder()
			l := s.echo.NewContext(request, w)
			l.Set("userID", "token")

			err := s.h.DeleteAllURLsByUserID(l)
			assert.NoError(t, err)

			assert.Equal(t, test.expectedCode, w.Code)
			s.ctrl.Finish()
		})
	}
}

func (s *URLHandlerTestSuite) TestDeleteAllURLsByUserID_BadRequest() {
	testCases := []struct {
		name         string
		method       string
		body         string
		expectedCode int
		path         string
		requestBody  string
		expectedBody string
	}{
		{
			name:         "BadRequest - empty request",
			method:       http.MethodDelete,
			expectedCode: http.StatusBadRequest,
			path:         "http://localhost:8080/api/user/urls",
			requestBody:  "",
			expectedBody: "",
		},
	}

	for _, test := range testCases {
		s.T().Run(test.name, func(t *testing.T) {
			s.urlService.EXPECT().DeleteURLByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			r, jsonErr := json.Marshal(test.requestBody)
			require.NoError(t, jsonErr)
			request := httptest.NewRequest(test.method, test.path, strings.NewReader(string(r)))
			request.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			l := s.echo.NewContext(request, w)
			l.Set("userID", "token")

			err := s.h.DeleteAllURLsByUserID(l)
			assert.NoError(t, err)

			assert.Equal(t, test.expectedCode, w.Code)
			s.ctrl.Finish()
		})
	}
}

func (s *URLHandlerTestSuite) TestDeleteAllURLsByUserID_Success() {
	testCases := []struct {
		name         string
		method       string
		body         string
		expectedCode int
		path         string
		requestBody  []string
		expectedBody string
	}{
		{
			name:         "Success",
			method:       http.MethodDelete,
			expectedCode: http.StatusAccepted,
			path:         "http://localhost:8080/api/user/urls",
			requestBody:  []string{"NjQyYTU", "OWUyMzI", "ZjQwMWN"},
			expectedBody: "",
		},
	}

	for _, test := range testCases {
		test := test
		s.T().Run(test.name, func(t *testing.T) {
			s.urlService.EXPECT().DeleteURLByUserID(gomock.Any(), gomock.Any(), gomock.Any()).MaxTimes(3).Return(nil)
			r, jsonErr := json.Marshal(test.requestBody)
			require.NoError(t, jsonErr)
			request := httptest.NewRequest(test.method, test.path, strings.NewReader(string(r)))
			request.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			l := s.echo.NewContext(request, w)
			l.Set("userID", "token")

			err := s.h.DeleteAllURLsByUserID(l)
			assert.NoError(t, err)

			assert.Equal(t, test.expectedCode, w.Code)
			s.ctrl.Finish()
		})
	}
}

func (s *URLHandlerTestSuite) TestFindAllURLByUserID_Unauthorized() {
	defer func(echo *echo.Echo) {
		err := echo.Close()
		assert.NoError(s.T(), err)
	}(s.echo)

	testCases := []struct {
		name         string
		method       string
		body         string
		expectedCode int
		path         string
		expectedBody string
	}{
		{
			name:         "BadRequest - unauthorized",
			method:       http.MethodGet,
			expectedCode: http.StatusUnauthorized,
			path:         "http://localhost:8080/api/user/urls",
			expectedBody: "",
		},
	}

	for _, test := range testCases {
		s.T().Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.path, strings.NewReader(""))
			w := httptest.NewRecorder()
			s.echo.ServeHTTP(w, request)
			assert.Equal(t, test.expectedCode, w.Code)
		})
	}
}

func (s *URLHandlerTestSuite) TestFindAllURLByUserID_NoContent() {
	testCases := []struct {
		name         string
		method       string
		body         string
		expectedCode int
		path         string
		expectedBody string
	}{
		{
			name:         "No Content",
			method:       http.MethodGet,
			expectedCode: http.StatusNoContent,
			path:         "http://localhost:8080/api/user/urls",
			expectedBody: "",
		},
	}

	for _, test := range testCases {
		s.T().Run(test.name, func(t *testing.T) {
			s.urlService.EXPECT().GetAllByUserID(gomock.Any(), gomock.Any()).Times(1).Return(nil, urlErr.ErrURLNotFound)
			request := httptest.NewRequest(test.method, test.path, strings.NewReader(""))
			w := httptest.NewRecorder()
			l := s.echo.NewContext(request, w)
			l.Set("userID", "token")

			err := s.h.FindAllURLByUserID(l)
			require.NoError(t, err)

			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, w.Body.String())
			s.ctrl.Finish()
		})
	}
}

func (s *URLHandlerTestSuite) TestAddBatch_WrongMediaType() {
	testCases := []struct {
		name         string
		method       string
		body         string
		expectedCode int
		path         string
		expectedBody string
	}{
		{
			name:         "BadRequest - unsupported media type",
			method:       http.MethodPost,
			expectedCode: http.StatusUnsupportedMediaType,
			path:         "http://localhost:8080/api/shorten/batch",
			expectedBody: "Content-Type header is not application/json",
		},
	}

	for _, test := range testCases {
		s.T().Run(test.name, func(t *testing.T) {
			s.urlService.EXPECT().AddAll(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			request := httptest.NewRequest(test.method, test.path, strings.NewReader(""))
			w := httptest.NewRecorder()
			l := s.echo.NewContext(request, w)
			l.Set("userID", "token")

			err := s.h.AddBatch(l)
			require.NoError(t, err)

			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, w.Body.String())
			s.ctrl.Finish()
		})
	}
}

func (s *URLHandlerTestSuite) TestAddBatch_EmptyRequest() {
	testCases := []struct {
		name         string
		method       string
		body         []string
		expectedCode int
		path         string
		expectedBody string
	}{
		{
			name:         "BadRequest - empty request",
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			path:         "http://localhost:8080/api/shorten/batch",
			body:         []string{},
			expectedBody: "Error: empty batch request",
		},
	}

	for _, test := range testCases {
		s.T().Run(test.name, func(t *testing.T) {
			s.urlService.EXPECT().AddAll(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			body, jsonErr := json.Marshal(test.body)
			require.NoError(t, jsonErr)
			request := httptest.NewRequest(test.method, test.path, strings.NewReader(string(body)))
			w := httptest.NewRecorder()
			l := s.echo.NewContext(request, w)
			request.Header.Set("Content-Type", "application/json")
			l.Set("userID", "token")

			err := s.h.AddBatch(l)
			require.NoError(t, err)

			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, w.Body.String())
			s.ctrl.Finish()
		})
	}
}

func (s *URLHandlerTestSuite) TestAddBatch_Success() {
	shortURL1 := dto.URLBatchResponse{
		CorrelationID: "1",
		ShortenedURL:  "1",
	}
	shortURL2 := dto.URLBatchResponse{
		CorrelationID: "2",
		ShortenedURL:  "2",
	}
	shortURLs := []dto.URLBatchResponse{shortURL1, shortURL2}

	fullURL1 := dto.URLBatchRequest{
		CorrelationID: "1",
		OriginalURL:   "1",
	}
	fullURL2 := dto.URLBatchRequest{
		CorrelationID: "2",
		OriginalURL:   "2",
	}
	request := []dto.URLBatchRequest{fullURL1, fullURL2}

	testCases := []struct {
		name         string
		method       string
		body         []dto.URLBatchRequest
		expectedCode int
		path         string
		expectedBody []dto.URLBatchResponse
	}{
		{
			name:         "Success",
			method:       http.MethodPost,
			expectedCode: http.StatusCreated,
			path:         "http://localhost:8080/api/shorten/batch",
			body:         request,
			expectedBody: shortURLs,
		},
	}

	for _, test := range testCases {
		s.T().Run(test.name, func(t *testing.T) {
			s.urlService.EXPECT().AddAll(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(shortURLs, nil)
			body, jsonErr := json.Marshal(test.body)
			require.NoError(t, jsonErr)
			req := httptest.NewRequest(test.method, test.path, strings.NewReader(string(body)))
			w := httptest.NewRecorder()
			l := s.echo.NewContext(req, w)
			req.Header.Set("Content-Type", "application/json")
			l.Set("userID", "token")

			err := s.h.AddBatch(l)
			require.NoError(t, err)

			assert.Equal(t, test.expectedCode, w.Code)
			var result []dto.URLBatchResponse
			jsonErr = json.Unmarshal(w.Body.Bytes(), &result)
			require.NoError(t, jsonErr)
			assert.Equal(t, test.expectedBody, result)
			assert.Equal(t, len(test.expectedBody), len(result))
			s.ctrl.Finish()
		})
	}
}

func (s *URLHandlerTestSuite) TestAddShorten_EmptyRequest() {
	testCases := []struct {
		name         string
		method       string
		body         dto.URLRequest
		expectedCode int
		path         string
		expectedBody string
	}{
		{
			name:         "BadRequest - empty request",
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			path:         "http://localhost:8080/api/shorten",
			body:         dto.URLRequest{},
			expectedBody: "Error: Unable to handle empty request",
		},
	}

	for _, test := range testCases {
		s.T().Run(test.name, func(t *testing.T) {
			s.urlService.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			body, jsonErr := json.Marshal(test.body)
			require.NoError(t, jsonErr)
			request := httptest.NewRequest(test.method, test.path, strings.NewReader(string(body)))
			w := httptest.NewRecorder()
			l := s.echo.NewContext(request, w)
			request.Header.Set("Content-Type", "application/json")
			l.Set("userID", "token")

			err := s.h.AddShorten(l)
			require.NoError(t, err)

			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, w.Body.String())
			s.ctrl.Finish()
		})
	}
}

func (s *URLHandlerTestSuite) TestAddShorten_WrongMediaType() {
	testCases := []struct {
		name         string
		method       string
		body         string
		expectedCode int
		path         string
		expectedBody string
	}{
		{
			name:         "BadRequest - unsupported media type",
			method:       http.MethodPost,
			expectedCode: http.StatusUnsupportedMediaType,
			path:         "http://localhost:8080/api/shorten",
			expectedBody: "Content-Type header is not application/json",
		},
	}

	for _, test := range testCases {
		s.T().Run(test.name, func(t *testing.T) {
			s.urlService.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			request := httptest.NewRequest(test.method, test.path, strings.NewReader(""))
			w := httptest.NewRecorder()
			l := s.echo.NewContext(request, w)
			l.Set("userID", "token")

			err := s.h.AddShorten(l)
			require.NoError(t, err)

			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, w.Body.String())
			s.ctrl.Finish()
		})
	}
}

func (s *URLHandlerTestSuite) TestAddShorten_InternalServerError() {
	requestBody := dto.URLRequest{URL: "https://example.com"}
	respErr := errors.New("internal server error")

	testCases := []struct {
		name         string
		method       string
		body         dto.URLRequest
		expectedCode int
		path         string
		expectedBody string
	}{
		{
			name:         "InternalServerError",
			method:       http.MethodPost,
			expectedCode: http.StatusInternalServerError,
			path:         "http://localhost:8080/api/shorten",
			body:         requestBody,
			expectedBody: "Unknown error: internal server error",
		},
	}

	for _, test := range testCases {
		s.T().Run(test.name, func(t *testing.T) {
			s.urlService.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil, respErr)
			body, jsonErr := json.Marshal(test.body)
			require.NoError(t, jsonErr)
			request := httptest.NewRequest(test.method, test.path, strings.NewReader(string(body)))
			w := httptest.NewRecorder()
			l := s.echo.NewContext(request, w)
			request.Header.Set("Content-Type", "application/json")
			l.Set("userID", "token")
			l.Set("userID", "token")

			err := s.h.AddShorten(l)
			require.NoError(t, err)

			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, w.Body.String())
			s.ctrl.Finish()
		})
	}
}

func (s *URLHandlerTestSuite) TestAddShorten_Success() {
	requestBody := dto.URLRequest{URL: "https://example.com"}
	url := &model.URL{Original: "https://example.com", Shortened: "test"}
	responseBody := dto.URLResponse{Result: url.Shortened}

	testCases := []struct {
		name         string
		method       string
		body         dto.URLRequest
		expectedCode int
		path         string
		expectedBody dto.URLResponse
	}{
		{
			name:         "Success",
			method:       http.MethodPost,
			expectedCode: http.StatusCreated,
			path:         "http://localhost:8080/api/shorten",
			body:         requestBody,
			expectedBody: responseBody,
		},
	}

	for _, test := range testCases {
		s.T().Run(test.name, func(t *testing.T) {
			s.urlService.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(url, nil)
			body, jsonErr := json.Marshal(test.body)
			require.NoError(t, jsonErr)
			request := httptest.NewRequest(test.method, test.path, strings.NewReader(string(body)))
			w := httptest.NewRecorder()
			l := s.echo.NewContext(request, w)
			request.Header.Set("Content-Type", "application/json")
			l.Set("userID", "token")
			l.Set("userID", "token")

			err := s.h.AddShorten(l)
			require.NoError(t, err)

			var result dto.URLResponse
			jsonErr = json.Unmarshal(w.Body.Bytes(), &result)
			require.NoError(t, jsonErr)
			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, result)
			s.ctrl.Finish()
		})
	}
}

func (s *URLHandlerTestSuite) TestAddShorten_UrlAlreadyExists() {
	requestBody := dto.URLRequest{URL: "https://example.com"}
	url := &model.URL{Original: "https://example.com", Shortened: "test"}
	responseBody := dto.URLResponse{Result: url.Shortened}
	mockErr := urlErr.ErrURLAlreadyExists

	testCases := []struct {
		name         string
		method       string
		body         dto.URLRequest
		expectedCode int
		path         string
		expectedBody dto.URLResponse
	}{
		{
			name:         "Status conflict - url already exists",
			method:       http.MethodPost,
			expectedCode: http.StatusConflict,
			path:         "http://localhost:8080/api/shorten",
			body:         requestBody,
			expectedBody: responseBody,
		},
	}

	for _, test := range testCases {
		s.T().Run(test.name, func(t *testing.T) {
			s.urlService.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(url, mockErr)
			body, jsonErr := json.Marshal(test.body)
			require.NoError(t, jsonErr)
			request := httptest.NewRequest(test.method, test.path, strings.NewReader(string(body)))
			w := httptest.NewRecorder()
			l := s.echo.NewContext(request, w)
			request.Header.Set("Content-Type", "application/json")
			l.Set("userID", "token")
			l.Set("userID", "token")

			err := s.h.AddShorten(l)
			require.NoError(t, err)

			var result dto.URLResponse
			jsonErr = json.Unmarshal(w.Body.Bytes(), &result)
			require.NoError(t, jsonErr)
			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, result)
			s.ctrl.Finish()
		})
	}
}

func (s *URLHandlerTestSuite) TestAddURL_EmptyRequest() {
	testCases := []struct {
		name         string
		method       string
		body         string
		expectedCode int
		path         string
		expectedBody string
	}{
		{
			name:         "BadRequest - empty request",
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			path:         "http://localhost:8080/",
			body:         "",
			expectedBody: "Error: Unable to handle empty request",
		},
	}

	for _, test := range testCases {
		s.T().Run(test.name, func(t *testing.T) {
			s.urlService.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			request := httptest.NewRequest(test.method, test.path, strings.NewReader(test.body))
			w := httptest.NewRecorder()
			l := s.echo.NewContext(request, w)
			request.Header.Set("Content-Type", "application/json")
			l.Set("userID", "token")

			err := s.h.AddURL(l)
			require.NoError(t, err)

			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, w.Body.String())
			s.ctrl.Finish()
		})
	}
}

func (s *URLHandlerTestSuite) TestAddURL_InternalServerError() {
	requestBody := "https://example.com"
	respErr := errors.New("internal server error")

	testCases := []struct {
		name         string
		method       string
		body         string
		expectedCode int
		path         string
		expectedBody string
	}{
		{
			name:         "InternalServerError",
			method:       http.MethodPost,
			expectedCode: http.StatusInternalServerError,
			path:         "http://localhost:8080/api/shorten",
			body:         requestBody,
			expectedBody: "Unknown error: internal server error",
		},
	}

	for _, test := range testCases {
		s.T().Run(test.name, func(t *testing.T) {
			s.urlService.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil, respErr)
			request := httptest.NewRequest(test.method, test.path, strings.NewReader(test.body))
			w := httptest.NewRecorder()
			l := s.echo.NewContext(request, w)
			request.Header.Set("Content-Type", "application/json")
			l.Set("userID", "token")
			l.Set("userID", "token")

			err := s.h.AddURL(l)
			require.NoError(t, err)

			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, w.Body.String())
			s.ctrl.Finish()
		})
	}
}

func (s *URLHandlerTestSuite) TestAddURL_Success() {
	requestBody := "https://example.com"
	url := &model.URL{Original: "https://example.com", Shortened: "test"}
	responseBody := url.Shortened

	testCases := []struct {
		name         string
		method       string
		body         string
		expectedCode int
		path         string
		expectedBody string
	}{
		{
			name:         "Success",
			method:       http.MethodPost,
			expectedCode: http.StatusCreated,
			path:         "http://localhost:8080/",
			body:         requestBody,
			expectedBody: responseBody,
		},
	}

	for _, test := range testCases {
		s.T().Run(test.name, func(t *testing.T) {
			s.urlService.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(url, nil)
			request := httptest.NewRequest(test.method, test.path, strings.NewReader(test.body))
			w := httptest.NewRecorder()
			l := s.echo.NewContext(request, w)
			request.Header.Set("Content-Type", "application/json")
			l.Set("userID", "token")
			l.Set("userID", "token")

			err := s.h.AddURL(l)
			require.NoError(t, err)

			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, w.Body.String())
			s.ctrl.Finish()
		})
	}
}

func (s *URLHandlerTestSuite) TestAddURL_UrlAlreadyExists() {
	requestBody := "https://example.com"
	url := &model.URL{Original: "https://example.com", Shortened: "test"}
	responseBody := url.Shortened
	mockErr := urlErr.ErrURLAlreadyExists

	testCases := []struct {
		name         string
		method       string
		body         string
		expectedCode int
		path         string
		expectedBody string
	}{
		{
			name:         "Status conflict - URL already exists",
			method:       http.MethodPost,
			expectedCode: http.StatusConflict,
			path:         "http://localhost:8080/",
			body:         requestBody,
			expectedBody: responseBody,
		},
	}

	for _, test := range testCases {
		s.T().Run(test.name, func(t *testing.T) {
			s.urlService.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(url, mockErr)
			request := httptest.NewRequest(test.method, test.path, strings.NewReader(test.body))
			w := httptest.NewRecorder()
			l := s.echo.NewContext(request, w)
			request.Header.Set("Content-Type", "application/json")
			l.Set("userID", "token")
			l.Set("userID", "token")

			err := s.h.AddURL(l)
			require.NoError(t, err)

			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, w.Body.String())
			s.ctrl.Finish()
		})
	}
}

func (s *URLHandlerTestSuite) TestClearALL_InternalServerError() {
	testCases := []struct {
		name         string
		method       string
		body         string
		expectedCode int
		path         string
		expectedBody string
	}{
		{
			name:         "InternalServerError",
			method:       http.MethodDelete,
			expectedCode: http.StatusInternalServerError,
			path:         "http://localhost:8080/",
			body:         "",
			expectedBody: "Unknown error: internal server error",
		},
	}

	for _, test := range testCases {
		s.T().Run(test.name, func(t *testing.T) {
			s.urlService.EXPECT().DeleteAll(gomock.Any()).Times(1).Return(errors.New("internal server error"))
			request := httptest.NewRequest(test.method, test.path, strings.NewReader(test.body))
			w := httptest.NewRecorder()
			l := s.echo.NewContext(request, w)
			l.Set("userID", "token")

			err := s.h.ClearAll(l)
			require.NoError(t, err)

			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, w.Body.String())
			s.ctrl.Finish()
		})
	}
}

func (s *URLHandlerTestSuite) TestClearALL_Success() {
	testCases := []struct {
		name         string
		method       string
		body         string
		expectedCode int
		path         string
		expectedBody string
	}{
		{
			name:         "Success",
			method:       http.MethodDelete,
			expectedCode: http.StatusOK,
			path:         "http://localhost:8080/",
			body:         "",
			expectedBody: "All data deleted",
		},
	}

	for _, test := range testCases {
		s.T().Run(test.name, func(t *testing.T) {
			s.urlService.EXPECT().DeleteAll(gomock.Any()).Times(1).Return(nil)
			request := httptest.NewRequest(test.method, test.path, strings.NewReader(test.body))
			w := httptest.NewRecorder()
			l := s.echo.NewContext(request, w)
			l.Set("userID", "token")

			err := s.h.ClearAll(l)
			require.NoError(t, err)

			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, w.Body.String())
			s.ctrl.Finish()
		})
	}
}

func (s *URLHandlerTestSuite) TestFindURL_Success() {
	url := &model.URL{Original: "https://example.com", Shortened: "test"}
	responseBody := url.Original

	testCases := []struct {
		name         string
		method       string
		body         string
		expectedCode int
		path         string
		expectedBody string
	}{
		{
			name:         "Success",
			method:       http.MethodPost,
			expectedCode: http.StatusTemporaryRedirect,
			path:         "http://localhost:8080/test",
			expectedBody: "",
		},
	}

	for _, test := range testCases {
		s.T().Run(test.name, func(t *testing.T) {
			s.urlService.EXPECT().GetByyID(gomock.Any(), url.Shortened).Times(1).Return(responseBody, nil)
			request := httptest.NewRequest(test.method, test.path, strings.NewReader(test.body))
			w := httptest.NewRecorder()
			l := s.echo.NewContext(request, w)
			request.Header.Set("Content-Type", "application/json")
			l.Set("userID", "token")
			l.Set("userID", "token")

			err := s.h.FindURL(l)
			require.NoError(t, err)

			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, w.Body.String())
			s.ctrl.Finish()
		})
	}
}

func (s *URLHandlerTestSuite) TestFindURL_EmptyRequest() {
	testCases := []struct {
		name         string
		method       string
		body         string
		expectedCode int
		path         string
		expectedBody string
	}{
		{
			name:         "BadRequest - empty request",
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			path:         "http://localhost:8080/",
			expectedBody: "Error: Unable to handle empty request",
		},
	}

	for _, test := range testCases {
		s.T().Run(test.name, func(t *testing.T) {
			s.urlService.EXPECT().GetByyID(gomock.Any(), gomock.Any()).Times(0)
			request := httptest.NewRequest(test.method, test.path, strings.NewReader(test.body))
			w := httptest.NewRecorder()
			l := s.echo.NewContext(request, w)
			request.Header.Set("Content-Type", "application/json")
			l.Set("userID", "token")
			l.Set("userID", "token")

			err := s.h.FindURL(l)
			require.NoError(t, err)

			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, w.Body.String())
			defer s.ctrl.Finish()
		})
	}
}
