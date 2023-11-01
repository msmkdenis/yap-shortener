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
	"github.com/msmkdenis/yap-shortener/internal/utils"
	"go.uber.org/zap"
)

const (
	perm         = 0755
)

type FileURLRepository struct {
	mu          sync.RWMutex
	fileStorage *os.File
	logger      *zap.Logger
}

func NewFileURLRepository(path string, logger *zap.Logger) *FileURLRepository {
	path = filepath.FromSlash(path)

	dir := filepath.Dir(path)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		logger.Info(fmt.Sprintf("Creating directory: %s", dir))
		err = os.Mkdir(dir, perm)
		if err != nil {
			logger.Fatal(fmt.Sprintf("Error while creating directory: %s", dir), zap.Error(err))
		}
		logger.Info(fmt.Sprintf("Directory %s was created", dir))
	}

	logger.Info(fmt.Sprintf("Creating file: %s", path))
	file, err := os.OpenFile(path, os.O_CREATE, perm)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Unable to create file: %s", path), zap.Error(err))
	}
	logger.Info(fmt.Sprintf("FileStorage %s was created", file.Name()))
	defer file.Close()

	return &FileURLRepository{
		fileStorage: file,
		logger:      logger,
		mu:          sync.RWMutex{},
	}
}

func (r *FileURLRepository) Insert(url model.URL) (*model.URL, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.logger.Info(fmt.Sprintf("Opening file: %s", r.fileStorage.Name()))
	file, err := os.OpenFile(r.fileStorage.Name(), os.O_RDWR|os.O_APPEND, perm)
	if err != nil {
		return nil, apperrors.NewValueError("unable to open file", utils.Caller(), err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(url)
	if err != nil {
		return nil, apperrors.NewValueError("unable to encode to file", utils.Caller(), err)
	}

	return &url, nil
}

func (r *FileURLRepository) SelectByID(key string) (*model.URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	file, err := os.OpenFile(r.fileStorage.Name(), os.O_RDONLY, perm)
	if err != nil {
		return nil, apperrors.NewValueError("unable to open file", utils.Caller(), err)
	}

	decoder := json.NewDecoder(file)
	defer file.Close()

	var url model.URL
	for {
		err := decoder.Decode(&url)
		if err == io.EOF {
			return nil, apperrors.NewValueError(fmt.Sprintf("Url with id %s not found", key), utils.Caller(), apperrors.ErrorURLNotFound)
		}
		if err != nil {
			return nil, apperrors.NewValueError("unable to decode from file", utils.Caller(), err)
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

	file, err := os.OpenFile(r.fileStorage.Name(), os.O_RDONLY, perm)
	if err != nil {
		return nil, apperrors.NewValueError("unable to open file", utils.Caller(), err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	urls := make([]model.URL, 0)

	for {
		var url model.URL
		err := decoder.Decode(&url)
		if err == io.EOF {
			r.logger.Info("Reached end of file while decoding", zap.Error(err))
			break
		}
		if err != nil {
			return nil, apperrors.NewValueError("unable to decode from file", utils.Caller(), err)
		}
		urls = append(urls, url)
	}

	return urls, nil
}

func (r *FileURLRepository) DeleteAll() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := os.Truncate(r.fileStorage.Name(), 0); err != nil {
		return apperrors.NewValueError(fmt.Sprintf("Failed to truncate file: %s", r.fileStorage.Name()), utils.Caller(), err)
	}
	return nil
}
