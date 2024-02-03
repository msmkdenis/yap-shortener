// Package middleware various middleware.
package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/msmkdenis/yap-shortener/pkg/compressor"
)

// Decompress returns a middleware that decompresses the request body if it is encoded with compressor.
func Decompress() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Header.Get("Content-Encoding") != "gzip" {
				return next(c)
			}
			b := c.Request().Body
			reader, err := compressor.NewReader(b, "gzip")
			//decompressingReader, err := compressor.NewGzipReader(b)
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
			if strings.Contains(c.Request().Header.Get("Accept-Encoding"), "gzip") {
				rw := c.Response().Writer
				cw := compressor.NewWriter(rw, "gzip")
				cw.Reset(rw)
			}
			return next(c)
		}
	}
}
