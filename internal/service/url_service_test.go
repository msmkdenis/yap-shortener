package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	mock "github.com/msmkdenis/yap-shortener/internal/mocks"
	"github.com/msmkdenis/yap-shortener/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"math/rand"
	"testing"
	"time"
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

func (u *URLServiceTestSuite) TestGetAll() {
	rnd := rand.NewSource(time.Now().Unix())
	data := make([]model.URL, 0, 10000)
	original := make([]string, 0, 10000)
	for i := 0; i < 10000; i++ {
		data = append(data, generateURL(rnd))
		original = append(original, data[i].Original)
	}

	testCases := []struct {
		name          string
		prepare       func()
		expectedCode  int
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
	}
	for _, test := range testCases {
		u.T().Run(test.name, func(t *testing.T) {
			if test.prepare != nil {
				test.prepare()
			}

			urls, err := u.urlService.GetAll(context.Background())
			assert.Equal(t, test.expectedError, err)
			assert.Equal(t, test.expectedBody, urls)
		})
	}
}

func BenchmarkURLUseCase_GetAll(b *testing.B) {
	s := new(URLServiceTestSuite)
	s.SetT(&testing.T{})
	s.SetupSuite()

	rnd := rand.NewSource(time.Now().Unix())
	data := make([]model.URL, 0, 10000)
	original := make([]string, 0, 10000)
	for i := 0; i < 10000; i++ {
		data = append(data, generateURL(rnd))
		original = append(original, data[i].Original)
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
