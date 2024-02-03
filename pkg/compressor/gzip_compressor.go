package compressor

import (
	"compress/gzip"
	"io"
	"net/http"
)

type GzipWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func NewGzipWriter(w http.ResponseWriter) *GzipWriter {
	return &GzipWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *GzipWriter) Reset(rw http.ResponseWriter) {
	c.zw.Reset(rw)
}

// Header implements http.ResponseWriter.
func (c *GzipWriter) Header() http.Header {
	return c.w.Header()
}

// Write implements io.Writer.
func (c *GzipWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

// WriteHeader implements http.ResponseWriter.
func (c *GzipWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close implements io.Closer.
func (c *GzipWriter) Close() error {
	return c.zw.Close()
}

type GzipReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func NewGzipReader(r io.ReadCloser) (*GzipReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &GzipReader{
		r:  r,
		zr: zr,
	}, nil
}

// Read implements io.Reader.
func (c *GzipReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close implements io.Closer.
func (c *GzipReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
