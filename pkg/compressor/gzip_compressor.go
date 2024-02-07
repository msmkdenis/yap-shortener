package compressor

import (
	"compress/gzip"
	"io"
	"net/http"
)

// GzipWriter implements Writer with gzip compression
type GzipWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

// NewGzipWriter returns a new GzipWriter with the given http.ResponseWriter.
func NewGzipWriter(w http.ResponseWriter) *GzipWriter {
	return &GzipWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// Reset discards writer state with gzip compression
func (c *GzipWriter) Reset(rw http.ResponseWriter) {
	c.zw.Reset(rw)
}

// Header returns http.Header with gzip compression
func (c *GzipWriter) Header() http.Header {
	return c.w.Header()
}

// Write compresses data with gzip compression
func (c *GzipWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

// WriteHeader write header with gzip compression
func (c *GzipWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close writer with gzip compression
func (c *GzipWriter) Close() error {
	return c.zw.Close()
}

// GzipReader implements Reader with gzip decompression
type GzipReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

// NewGzipReader returns a new GzipReader with the given io.ReadCloser.
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

// Read implements io.Reader with gzip decompression
func (c *GzipReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close implements io.Closer with gzip decompression
func (c *GzipReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
