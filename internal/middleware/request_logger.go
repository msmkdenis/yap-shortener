package middleware

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type (
	RequestLogger struct {
		ReqLogger *zap.Logger
	}

	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func InitRequestLogger(logger *zap.Logger) *RequestLogger {
	l := &RequestLogger{
		ReqLogger: logger,
	}
	return l
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// RequestLogger logs each HTTP request.
func (r *RequestLogger) RequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			uri := c.Request().RequestURI

			method := c.Request().Method

			duration := time.Since(start)

			responseData := &responseData{}

			lw := loggingResponseWriter{
				ResponseWriter: c.Response().Writer,
				responseData:   responseData,
			}

			c.Response().Writer = &lw

			err := next(c)

			r.ReqLogger.Info("request_logger",
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
