package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/msmkdenis/yap-shortener/internal/dto"
	"math/rand"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	mock "github.com/msmkdenis/yap-shortener/internal/mocks"
	"github.com/msmkdenis/yap-shortener/internal/model"
	urlErr "github.com/msmkdenis/yap-shortener/internal/urlerr"
	"github.com/msmkdenis/yap-shortener/pkg/hasher"
)

type URLServiceTestSuite struct {
	suite.Suite
	logger        *zap.Logger
	urlRepository *mock.MockURLRepository
	urlService    *URLUseCase
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(URLServiceTestSuite))
}

func (u *URLServiceTestSuite) SetupSuite() {
	u.logger, _ = zap.NewProduction()
	u.urlRepository = mock.NewMockURLRepository(gomock.NewController(u.T()))
	u.urlService = NewURLService(u.urlRepository, u.logger)
}

func (u *URLServiceTestSuite) TestGetAllByUserId() {
	rnd := rand.NewSource(time.Now().Unix())
	data := make([]model.URL, 0, 100)
	batchResponse := make([]dto.URLBatchResponseByUserID, 0, 100)
	for i := 0; i < 100; i++ {
		data = append(data, generateURL(rnd))
		batchResponse = append(batchResponse, dto.URLBatchResponseByUserID{
			DeletedFlag: data[i].DeletedFlag,
			OriginalURL: data[i].Original,
			ShortURL:    data[i].Shortened,
		})
	}

	repoErr := errors.New("repository error")

	testCases := []struct {
		name          string
		prepare       func()
		expectedBody  []dto.URLBatchResponseByUserID
		expectedError error
	}{
		{
			name: "Successful return",
			prepare: func() {
				u.urlRepository.EXPECT().SelectAllByUserID(gomock.Any(), gomock.Any()).Return(data, nil)
			},
			expectedBody:  batchResponse,
			expectedError: nil,
		},
		{
			name: "Error return",
			prepare: func() {
				u.urlRepository.EXPECT().SelectAllByUserID(gomock.Any(), gomock.Any()).Return(nil, repoErr)
			},
			expectedBody: nil,
		},
	}
	for _, test := range testCases {
		u.T().Run(test.name, func(t *testing.T) {
			if test.prepare != nil {
				test.prepare()
			}

			urls, err := u.urlService.GetAllByUserID(context.Background(), uuid.New().String())
			assert.Equal(t, test.expectedBody, urls)
			if err != nil {
				assert.True(t, errors.Is(err, repoErr))
			} else {
				assert.Equal(t, test.expectedError, err)
			}
		})
	}
}

func (u *URLServiceTestSuite) TestDeleteURLByUserID() {
	repoErr := errors.New("repository error")

	testCases := []struct {
		name          string
		prepare       func()
		expectedError error
	}{
		{
			name: "Successful delete",
			prepare: func() {
				u.urlRepository.EXPECT().DeleteURLByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Error delete",
			prepare: func() {
				u.urlRepository.EXPECT().DeleteURLByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return(repoErr)
			},
		},
	}
	for _, test := range testCases {
		u.T().Run(test.name, func(t *testing.T) {
			if test.prepare != nil {
				test.prepare()
			}

			err := u.urlService.DeleteURLByUserID(context.Background(), uuid.New().String(), "shortened")
			switch test.name {
			case "Successful delete":
				assert.Equal(t, test.expectedError, err)
			case "Error delete":
				assert.True(t, errors.Is(err, repoErr))
			}
		})
	}
}

func (u *URLServiceTestSuite) TestAdd() {
	rnd := rand.NewSource(time.Now().Unix())
	s := generateString(10, rnd)
	host := generateString(4, rnd)
	userID := uuid.New().String()
	urlKey := hasher.GenerateMD5Hash(s)
	url := &model.URL{
		ID:          urlKey,
		Original:    s,
		Shortened:   host + "/" + urlKey,
		UserID:      userID,
		DeletedFlag: false,
	}

	repoErr := errors.New("repository error")

	testCases := []struct {
		name          string
		prepareSelect func()
		prepareInsert func()
		expectedBody  *model.URL
		expectedError error
	}{
		{
			name: "Successful add",
			prepareSelect: func() {
				u.urlRepository.EXPECT().SelectByID(gomock.Any(), urlKey).Return(nil, repoErr)
			},
			prepareInsert: func() {
				u.urlRepository.EXPECT().Insert(gomock.Any(), *url).Return(url, nil)
			},
			expectedBody:  url,
			expectedError: nil,
		},
		{
			name: "Successful return existing url",
			prepareSelect: func() {
				u.urlRepository.EXPECT().SelectByID(gomock.Any(), urlKey).Return(url, nil)
			},
			expectedBody:  url,
			expectedError: urlErr.ErrURLAlreadyExists,
		},
		{
			name: "Error while add",
			prepareSelect: func() {
				u.urlRepository.EXPECT().SelectByID(gomock.Any(), urlKey).Return(nil, repoErr)
			},
			prepareInsert: func() {
				u.urlRepository.EXPECT().Insert(gomock.Any(), *url).Return(nil, repoErr)
			},
			expectedBody:  nil,
			expectedError: repoErr,
		},
	}
	for _, test := range testCases {
		u.T().Run(test.name, func(t *testing.T) {
			if test.prepareSelect != nil {
				test.prepareSelect()
			}
			if test.prepareInsert != nil {
				test.prepareInsert()
			}

			savedURL, err := u.urlService.Add(context.Background(), s, host, userID)
			assert.Equal(t, test.expectedBody, savedURL)
			if err != nil {
				assert.True(t, errors.Is(err, test.expectedError))
			} else {
				assert.Equal(t, test.expectedError, err)
			}
		})
	}
}

func (u *URLServiceTestSuite) TestGetAll() {
	rnd := rand.NewSource(time.Now().Unix())
	data := make([]model.URL, 0, 100)
	original := make([]string, 0, 100)
	for i := 0; i < 100; i++ {
		data = append(data, generateURL(rnd))
		original = append(original, data[i].Original)
	}

	repoErr := errors.New("repository error")

	testCases := []struct {
		name          string
		prepare       func()
		expectedBody  []string
		expectedError error
	}{
		{
			name: "Successful return",
			prepare: func() {
				u.urlRepository.EXPECT().SelectAll(gomock.Any()).Return(data, nil)
			},
			expectedBody:  original,
			expectedError: nil,
		},
		{
			name: "Error return",
			prepare: func() {
				u.urlRepository.EXPECT().SelectAll(gomock.Any()).Return(nil, repoErr)
			},
			expectedBody: nil,
		},
	}
	for _, test := range testCases {
		u.T().Run(test.name, func(t *testing.T) {
			if test.prepare != nil {
				test.prepare()
			}

			urls, err := u.urlService.GetAll(context.Background())
			assert.Equal(t, test.expectedBody, urls)
			if err != nil {
				assert.True(t, errors.Is(err, repoErr))
			} else {
				assert.Equal(t, test.expectedError, err)
			}
		})
	}
}

func (u *URLServiceTestSuite) TestDeleteAll() {
	repoErr := errors.New("repository error")

	testCases := []struct {
		name          string
		prepare       func()
		expectedBody  []string
		expectedError error
	}{
		{
			name: "Successful delete",
			prepare: func() {
				u.urlRepository.EXPECT().DeleteAll(gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Error",
			prepare: func() {
				u.urlRepository.EXPECT().DeleteAll(gomock.Any()).Return(repoErr)
			},
			expectedError: repoErr,
		},
	}
	for _, test := range testCases {
		u.T().Run(test.name, func(t *testing.T) {
			if test.prepare != nil {
				test.prepare()
			}

			err := u.urlService.DeleteAll(context.Background())
			if err != nil {
				assert.True(t, errors.Is(err, repoErr))
			} else {
				assert.Equal(t, test.expectedError, err)
			}
		})
	}
}

func (u *URLServiceTestSuite) TestGetByID() {
	rnd := rand.NewSource(time.Now().Unix())
	s := generateString(10, rnd)
	host := generateString(4, rnd)
	userID := uuid.New().String()
	urlKey := hasher.GenerateMD5Hash(s)
	url := &model.URL{
		ID:          urlKey,
		Original:    s,
		Shortened:   host + "/" + urlKey,
		UserID:      userID,
		DeletedFlag: false,
	}

	repoErr := errors.New("repository error")

	testCases := []struct {
		name          string
		prepare       func()
		expectedBody  string
		expectedError error
	}{
		{
			name: "Successful get",
			prepare: func() {
				u.urlRepository.EXPECT().SelectByID(gomock.Any(), urlKey).Return(url, nil)
			},
			expectedBody:  url.Original,
			expectedError: nil,
		},
		{
			name: "Url deleted",
			prepare: func() {
				url.DeletedFlag = true
				u.urlRepository.EXPECT().SelectByID(gomock.Any(), urlKey).Return(url, nil)
			},
			expectedError: urlErr.ErrURLDeleted,
		},
		{
			name: "Error",
			prepare: func() {
				u.urlRepository.EXPECT().SelectByID(gomock.Any(), urlKey).Return(nil, repoErr)
			},
			expectedError: repoErr,
		},
	}
	for _, test := range testCases {
		u.T().Run(test.name, func(t *testing.T) {
			if test.prepare != nil {
				test.prepare()
			}

			url, err := u.urlService.GetByyID(context.Background(), urlKey)
			if err != nil {
				assert.True(t, errors.Is(err, test.expectedError))
			} else {
				assert.Equal(t, test.expectedBody, url)
			}
		})
	}
}

func (u *URLServiceTestSuite) TestAddAll() {
	rnd := rand.NewSource(time.Now().Unix())
	host := generateString(4, rnd)
	userID := uuid.New().String()

	urlBatchRequest := make([]dto.URLBatchRequest, 0, 20)
	urlBatchResponse := make([]dto.URLBatchResponse, 0, 20)
	urls := make([]model.URL, 0, 20)
	for i := 0; i < 20; i++ {
		s := generateString(10, rnd)
		request := dto.URLBatchRequest{
			CorrelationID: uuid.New().String(),
			OriginalURL:   s,
		}
		urlBatchRequest = append(urlBatchRequest, request)

		shortURL := hasher.GenerateMD5Hash(request.OriginalURL)
		url := model.URL{
			ID:            shortURL,
			Original:      request.OriginalURL,
			Shortened:     host + "/" + shortURL,
			CorrelationID: request.CorrelationID,
			UserID:        userID,
			DeletedFlag:   false,
		}
		urls = append(urls, url)

		response := dto.URLBatchResponse{
			CorrelationID: request.CorrelationID,
			ShortenedURL:  host + "/" + hasher.GenerateMD5Hash(request.OriginalURL),
		}
		urlBatchResponse = append(urlBatchResponse, response)
	}

	repoErr := errors.New("repository error")

	testCases := []struct {
		name          string
		prepare       func()
		expectedBody  []dto.URLBatchResponse
		expectedError error
	}{
		{
			name: "Successful add all",
			prepare: func() {
				u.urlRepository.EXPECT().InsertAllOrUpdate(gomock.Any(), gomock.Any()).Return(urls, nil)
			},
			expectedBody:  urlBatchResponse,
			expectedError: nil,
		},
		{
			name: "Error duplicated key",
			prepare: func() {
				urlBatchRequest[0].CorrelationID = uuid.New().String()
				urlBatchRequest[1].CorrelationID = urlBatchRequest[0].CorrelationID
			},
			expectedError: urlErr.ErrDuplicatedKeys,
		},
		{
			name: "Error",
			prepare: func() {
				urlBatchRequest[1].CorrelationID = uuid.New().String()
				u.urlRepository.EXPECT().InsertAllOrUpdate(gomock.Any(), gomock.Any()).Return(nil, repoErr)
			},
			expectedError: repoErr,
		},
	}
	for _, test := range testCases {
		u.T().Run(test.name, func(t *testing.T) {
			if test.prepare != nil {
				test.prepare()
			}

			savedURL, err := u.urlService.AddAll(context.Background(), urlBatchRequest, host, userID)
			fmt.Println(err)
			assert.Equal(t, test.expectedBody, savedURL)
			if err != nil {
				assert.True(t, errors.Is(err, test.expectedError))
			} else {
				assert.Equal(t, test.expectedError, err)
			}
		})
	}
}

func BenchmarkURLUseCase_GetAll(b *testing.B) {
	s := new(URLServiceTestSuite)
	s.SetT(&testing.T{})
	s.SetupSuite()

	rnd := rand.NewSource(time.Now().Unix())
	data := make([]model.URL, 0, 10000)
	for i := 0; i < 10000; i++ {
		data = append(data, generateURL(rnd))
	}
	b.StartTimer()
	b.Run("UrlServiceGetAll", func(b *testing.B) {
		for j := 0; j < b.N; j++ {
			s.urlRepository.EXPECT().SelectAll(gomock.Any()).Return(data, nil)
			_, _ = s.urlService.GetAll(context.Background())
		}
	})
	b.StopTimer()
}

func generateURL(rnd rand.Source) model.URL {
	url := model.URL{
		ID:            uuid.New().String(),
		Original:      generateString(10, rnd),
		Shortened:     generateString(5, rnd),
		CorrelationID: generateString(5, rnd),
		UserID:        uuid.New().String(),
		DeletedFlag:   false,
	}

	return url
}

func generateString(n int, generator rand.Source) string {
	result := make([]byte, 0, n)
	for i := 0; i < n; i++ {
		// генерируем случайное число
		randomNumber := generator.Int63()
		// английские буквы лежат в диапазоне от 97 до 122, поэтому:
		// 1) берем остаток от деления случайного числа на 26, получая диапазон [0,25]
		// 2) прибавляем к полученному числу 97 и получаем итоговый интервал: [97, 122].
		result = append(result, byte(randomNumber%26+97))
	}
	return string(result)
}
