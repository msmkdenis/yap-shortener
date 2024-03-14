package integrationurltests

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"

	"github.com/msmkdenis/yap-shortener/internal/api/httphandlers"
	"github.com/msmkdenis/yap-shortener/internal/config"
	"github.com/msmkdenis/yap-shortener/internal/middleware"
	"github.com/msmkdenis/yap-shortener/internal/repository/db"
	"github.com/msmkdenis/yap-shortener/internal/service"
	"github.com/msmkdenis/yap-shortener/pkg/jwtgen"
)

var cfgMock = &config.Config{
	TokenName:     "test",
	SecretKey:     "test",
	TrustedSubnet: "",
}

type IntegrationTestSuite struct {
	suite.Suite
	urlHandler    *httphandlers.URLShorten
	urlService    *service.URLUseCase
	urlRepository *db.PostgresURLRepository
	echo          *echo.Echo
	container     testcontainers.Container
	pool          *db.PostgresPool
	endpoint      string
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) SetupTest() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Error("Unable to initialize zap logger", zap.Error(err))
	}

	s.container, s.pool, err = setupTestDatabase()
	if err != nil {
		logger.Error("Unable to setup test database", zap.Error(err))
	}

	s.urlRepository = db.NewPostgresURLRepository(s.pool, logger)
	jwtManager := jwtgen.InitJWTManager(cfgMock.TokenName, cfgMock.SecretKey, logger)
	jwtCheckerCreator := middleware.InitJWTCheckerCreator(jwtManager, logger)
	jwtAuth := middleware.InitJWTAuth(jwtManager, logger)
	s.urlService = service.NewURLService(s.urlRepository, logger)
	s.echo = echo.New()
	s.endpoint, err = s.container.Endpoint(context.Background(), "httphandlers")
	if err != nil {
		logger.Error("Unable to get endpoint", zap.Error(err))
	}
	s.urlHandler = httphandlers.NewURLShorten(s.echo, s.urlService, s.endpoint, cfgMock.TrustedSubnet, jwtCheckerCreator, jwtAuth, logger, &sync.WaitGroup{})
}

func (s *IntegrationTestSuite) TestAddURL() {
	body := "https://example.com"

	testCases := []struct {
		name         string
		method       string
		header       http.Header
		prepare      func()
		path         string
		body         string
		expectedCode int
		expectedBody []byte
	}{
		{
			name:         "Success - 201",
			method:       http.MethodPost,
			path:         s.endpoint + "/",
			body:         body,
			expectedCode: http.StatusCreated,
			expectedBody: []byte(fmt.Sprintf("%s/Yzk4NGQ", s.endpoint)),
		},
		{
			name:         "Empty request - 400",
			method:       http.MethodPost,
			path:         s.endpoint + "/",
			body:         "",
			expectedCode: http.StatusBadRequest,
			expectedBody: []byte("Error: Unable to handle empty request"),
		},
		{
			name:         "UrlAlreadyExists - 409",
			method:       http.MethodPost,
			path:         s.endpoint + "/",
			body:         body,
			expectedCode: http.StatusConflict,
			expectedBody: []byte(fmt.Sprintf("%s/Yzk4NGQ", s.endpoint)),
		},
	}
	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, bytes.NewBuffer([]byte(tc.body)))

			req.Header = tc.header
			rec := httptest.NewRecorder()

			s.echo.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.Equal(t, tc.expectedBody, rec.Body.Bytes())
		})
	}
}

func (s *IntegrationTestSuite) TestAddURL_Context() {
	body := "https://example.com"

	testCases := []struct {
		name         string
		method       string
		header       http.Header
		prepare      func()
		path         string
		body         string
		expectedCode int
		expectedBody []byte
	}{
		{
			name:         "UnableToGetUserIDFromContext - 500",
			method:       http.MethodPost,
			path:         s.endpoint + "/",
			body:         body,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "Bad request - 400 (unable to read body)",
			method:       http.MethodPost,
			path:         s.endpoint + "/",
			body:         "",
			expectedCode: http.StatusBadRequest,
			expectedBody: []byte("Error: Unable to handle empty request"),
		},
	}
	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, bytes.NewBuffer([]byte(tc.body)))
			req.Header = tc.header
			rec := httptest.NewRecorder()

			c := s.echo.NewContext(req, rec)

			err := s.urlHandler.AddURL(c)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.Equal(t, tc.expectedBody, rec.Body.Bytes())
		})
	}
}

func (s *IntegrationTestSuite) TearDownTest() {
	defer s.container.Terminate(context.Background())
}

func setupTestDatabase() (testcontainers.Container, *db.PostgresPool, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Unable to initialize zap logger", err)
	}

	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_DB":       "yap-shortener-test",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_USER":     "postgres",
		},
	}
	dbContainer, err := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})
	if err != nil {
		return nil, nil, err
	}

	port, err := dbContainer.MappedPort(context.Background(), "5432")
	if err != nil {
		return nil, nil, err
	}
	host, err := dbContainer.Host(context.Background())
	if err != nil {
		return nil, nil, err
	}

	connection := fmt.Sprintf("user=postgres password=postgres host=%s database=yap-shortener-test sslmode=disable port=%d", host, port.Int())

	pool := initPostgresPool(connection, logger)

	return dbContainer, pool, err
}

func initPostgresPool(uri string, logger *zap.Logger) *db.PostgresPool {
	postgresPool, err := db.NewPostgresPool(uri, logger)
	if err != nil {
		logger.Fatal("Unable to connect to database", zap.Error(err))
	}

	migrations, err := db.NewMigrations(uri, logger)
	if err != nil {
		logger.Fatal("Unable to create migrations", zap.Error(err))
	}

	err = migrations.MigrateUp()
	if err != nil {
		logger.Fatal("Unable to up migrations", zap.Error(err))
	}

	logger.Info("Connected to database", zap.String("DSN", uri))
	return postgresPool
}
