package compressor

import (
	"errors"
	"io"
	"net/http"
)

type Writer interface {
	Reset(rw http.ResponseWriter)
	Header() http.Header
	Write(p []byte) (int, error)
	WriteHeader(statusCode int)
	Close() error
}

func NewWriter(w http.ResponseWriter, alg string) Writer {
	switch alg {
	case "gzip":
		return NewGzipWriter(w)
	}
	return nil
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
	return nil, errors.New("unknown compression algorithm")
}
