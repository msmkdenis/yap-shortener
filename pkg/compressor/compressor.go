package compressor

import (
	"errors"
	"io"
	"net/http"
)

var ErrUnknownCompressionAlgorithm = errors.New("url not found")

type Writer interface {
	Reset(rw http.ResponseWriter)
	Header() http.Header
	Write(p []byte) (int, error)
	WriteHeader(statusCode int)
	Close() error
}

func NewWriter(w http.ResponseWriter, alg string) (Writer, error) {
	switch alg {
	case "gzip":
		return NewGzipWriter(w), nil
	}
	return nil, ErrUnknownCompressionAlgorithm
}

type Reader interface {
	Read(p []byte) (n int, err error)
	Close() error
}

func NewReader(r io.ReadCloser, alg string) (Reader, error) {
	switch alg {
	case "gzip":
		reader, err := NewGzipReader(r)
		return reader, err
	}
	return nil, ErrUnknownCompressionAlgorithm
}
