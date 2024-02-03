// Package middleware various middleware.
package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/msmkdenis/yap-shortener/pkg/compressor"
	"net/http"
)

// Decompress returns a middleware that decompresses the request body if it is encoded with compressor.
func Decompress() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			decompress := c.Request().Header.Get("Content-Encoding")
			if decompress == "" {
				return next(c)
			}
			b := c.Request().Body
			reader, err := compressor.NewReader(b, decompress)
			if err == nil {
				c.Request().Body = reader
				return next(c)
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
}

// Compress returns a middleware function that compresses the response using compressor if the client supports it.
func Compress() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			compress := c.Response().Header().Get("Content-Encoding")
			if compress == "" {
				return next(c)
			}
			cw, err := compressor.NewWriter(c.Response().Writer, compress)
			if err != nil {
				cw.Reset(c.Response().Writer)
			}
			return next(c)
		}
	}
}
