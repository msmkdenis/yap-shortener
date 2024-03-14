// Package benchclient sends simple requests to the server.
package client

import (
	"math/rand"
	"os"
	"time"

	"github.com/msmkdenis/yap-shortener/internal/dto"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

// Config contains the configuration for the benchclient.
type Config struct {
	Address string `env:"CLIENT_ADDRESS" envDefault:"0.0.0.0:6060"`
}

// Run sends GET requests and exit after 30 seconds.
func Run() {
	time.AfterFunc(time.Second*30, func() {
		os.Exit(1)
	})

	logger, _ := zap.NewProduction()
	client := resty.New()

	rnd := rand.NewSource(time.Now().Unix())
	var sendRequest []dto.URLBatchRequest
	for i := 0; i < 10000; i++ {
		url := generateBatch(rnd)
		sendRequest = append(sendRequest, url)
	}

	_, err := client.R().
		SetBody(sendRequest).
		SetHeaders(map[string]string{"Content-Type": "application/json"}).
		Post("http://localhost:8080/api/shorten/batch")
	if err != nil {
		logger.Fatal("Error", zap.Error(err))
	}

	for i := 0; i < 10000; i++ {
		time.Sleep(100 * time.Millisecond)
		_, err := client.R().
			SetHeaders(map[string]string{"Content-Type": "application/json"}).
			Get("http://localhost:8080/")
		if err != nil {
			logger.Info("Error", zap.Error(err))
		}
	}
}

func generateBatch(rnd rand.Source) dto.URLBatchRequest {
	url := dto.URLBatchRequest{
		CorrelationID: generateString(10, rnd),
		OriginalURL:   generateString(10, rnd),
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
