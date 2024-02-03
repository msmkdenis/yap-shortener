package compressor

import (
	"errors"
	"io"
	"net/http"
)

// ErrUnknownCompressionAlgorithm error for unknown compression algorithm
var ErrUnknownCompressionAlgorithm = errors.New("url not found")

// Writer interface for compressing
type Writer interface {
	Reset(rw http.ResponseWriter)
	Header() http.Header
	Write(p []byte) (int, error)
	WriteHeader(statusCode int)
	Close() error
}

// NewWriter factory method returns a new writer for compressing
//
// Parameters: http.ResponseWriter, string - the compression algorithm
func NewWriter(w http.ResponseWriter, alg string) (Writer, error) {
	switch alg {
	case "gzip":
		return NewGzipWriter(w), nil
	}
	return nil, ErrUnknownCompressionAlgorithm
}

// Reader interface for decompressing
type Reader interface {
	Read(p []byte) (n int, err error)
	Close() error
}

// NewReader factory method returns a new reader for decompressing
//
// Parameters: io.ReadCloser, string - the decompression algorithm
func NewReader(r io.ReadCloser, alg string) (Reader, error) {
	switch alg {
	case "gzip":
		reader, err := NewGzipReader(r)
		return reader, err
	}
	return nil, ErrUnknownCompressionAlgorithm
}
