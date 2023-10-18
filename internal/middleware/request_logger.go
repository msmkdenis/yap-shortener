package middleware

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type RequestLogger struct {
	ReqLogger *zap.Logger
}

func InitRequestLogger(logger *zap.Logger) *RequestLogger {
	fmt.Println(logger)
	l := &RequestLogger{
		ReqLogger: logger,
	}
	return l
}

type (
	// берём структуру для хранения сведений об ответе
	responseData struct {
		status int
		size   int
	}

	// добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

func (r *RequestLogger) RequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			uri := c.Request().RequestURI

			method := c.Request().Method

			duration := time.Since(start)

			responseData := &responseData{
				status: 0,
				size:   0,
			}

			lw := loggingResponseWriter{
				ResponseWriter: c.Response().Writer, // встраиваем оригинальный http.ResponseWriter
				responseData:   responseData,
			}

			c.Response().Writer = &lw

			err := next(c)

			r.ReqLogger.Info("post_logger",
				zap.String("URI", uri),
				zap.String("method", method),
				zap.Duration("duration", duration),
				zap.Int("response_code", responseData.status),
				zap.Int("response_body_size", responseData.size),
			)
			return err
		}
	}
}
