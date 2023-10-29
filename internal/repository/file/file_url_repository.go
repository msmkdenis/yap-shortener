package file

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/msmkdenis/yap-shortener/internal/apperrors"
	"github.com/msmkdenis/yap-shortener/internal/model"
	"go.uber.org/zap"
)

const (
	openFileFlag = os.O_RDWR | os.O_CREATE | os.O_APPEND
	perm         = 0755
)

type FileURLRepository struct {
	mu              sync.RWMutex
	fileStoragePath string
	logger          *zap.Logger
}

func NewFileURLRepository(path string, logger *zap.Logger) *FileURLRepository {
	path = filepath.FromSlash(path)

	dir := filepath.Dir(path)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, perm)
		if err != nil {
			logger.Fatal(fmt.Sprintf("Error while creating directory: %s", dir), zap.Error(err))
			return nil
		}
	}

	return &FileURLRepository{
		fileStoragePath: path,
		logger:          logger,
		mu:              sync.RWMutex{},
	}
}

func (r *FileURLRepository) Insert(url model.URL) (*model.URL, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.logger.Info(fmt.Sprintf("Opening file: %s", r.fileStoragePath))
	file, err := os.OpenFile(r.fileStoragePath, openFileFlag, perm)
	if err != nil {
		r.logger.Fatal(fmt.Sprintf("Unable to open file: %s", r.fileStoragePath), zap.Error(err))
		return nil, err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(url)


	if err != nil {
		r.logger.Error("Could not encode url to file", zap.Error(err))
		return nil, err
	}

	return &url, nil
}

func (r *FileURLRepository) SelectByID(key string) (*model.URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	file, err := os.OpenFile(r.fileStoragePath, openFileFlag, perm)
	if err != nil {
		r.logger.Fatal("Unable to open file", zap.Error(err))
		return nil, err
	}

	decoder := json.NewDecoder(file)
	defer file.Close()

	var url model.URL
	for {
		err := decoder.Decode(&url)
		if err == io.EOF {
			r.logger.Debug(fmt.Sprintf("url with id %s not found", key), zap.Error(err))
			return nil, apperrors.ErrorUrlNotFound
		}
		if err != nil {
			r.logger.Error("Could not decode url from file", zap.Error(err))
			return nil, err
		}
		if url.ID == key {
			break
		}
	}


	return &url, nil
}

func (r *FileURLRepository) SelectAll() ([]model.URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	file, err := os.OpenFile(r.fileStoragePath, openFileFlag, perm)
	if err != nil {
		r.logger.Fatal("Unable to open file", zap.Error(err))
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	urls := make([]model.URL, 0)

	for {
		var url model.URL
		err := decoder.Decode(&url)
		if err == io.EOF {
			r.logger.Debug("Reached end of file while decoding", zap.Error(err))
			break
		}
		if err != nil {
			r.logger.Error("Could not decode url from file", zap.Error(err))
			return nil, err
		}
		urls = append(urls, url)
	}

	return urls, nil
}

func (r *FileURLRepository) DeleteAll() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := os.Truncate(r.fileStoragePath, 0); err != nil {
		r.logger.Error(fmt.Sprintf("Failed to truncate file: %s", r.fileStoragePath), zap.Error(err))
	}
}
