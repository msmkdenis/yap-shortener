// Package file contains the file repository implementation.
package file

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"

	"github.com/msmkdenis/yap-shortener/internal/model"
	urlErr "github.com/msmkdenis/yap-shortener/internal/url_err"
	"github.com/msmkdenis/yap-shortener/pkg/apperr"
)

const (
	perm = 0o755
)

// FileURLRepository represents a file-based implementation of the URLRepository interface.
type URLRepository struct {
	mu          sync.RWMutex
	fileStorage *os.File
	logger      *zap.Logger
}

// NewFileURLRepository creates a new URLRepository from the given path and logger.
// Tries to create the directory if it doesn't exist.
// Tries to create the file if it doesn't exist.
func NewFileURLRepository(path string, logger *zap.Logger) (*URLRepository, error) {
	path = filepath.FromSlash(path)

	dir := filepath.Dir(path)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		logger.Info(fmt.Sprintf("Creating directory: %s", dir))
		err = os.Mkdir(dir, perm)
		if err != nil {
			return nil, apperr.NewValueError(fmt.Sprintf("Unable to create directory: %s", dir), apperr.Caller(), err)
		}
		logger.Info(fmt.Sprintf("Directory %s was created", dir))
	}

	logger.Info(fmt.Sprintf("Creating file: %s", path))
	file, err := os.OpenFile(path, os.O_CREATE, perm)
	if err != nil {
		return nil, apperr.NewValueError(fmt.Sprintf("Unable to create file: %s", path), apperr.Caller(), err)
	}
	logger.Info(fmt.Sprintf("FileStorage %s was created", file.Name()))
	defer file.Close()

	return &URLRepository{
		fileStorage: file,
		logger:      logger,
		mu:          sync.RWMutex{},
	}, nil
}

// DeleteURLByUserID deletes a URL by user ID from the file
func (r *URLRepository) DeleteURLByUserID(ctx context.Context, userID string, shortURL string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.logger.Info(fmt.Sprintf("Opening file: %s", r.fileStorage.Name()))
	file, openFileErr := os.OpenFile(r.fileStorage.Name(), os.O_RDWR|os.O_APPEND, perm)
	if openFileErr != nil {
		return apperr.NewValueError("unable to open file", apperr.Caller(), openFileErr)
	}
	defer file.Close()

	// Read all urls from file, update with new flag, store urls to save (with updated ones) in urlsToSave slice
	decoder := json.NewDecoder(file)
	var urlsToSave []model.URL
	for {
		var existingURL model.URL
		err := decoder.Decode(&existingURL)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return apperr.NewValueError("unable to decode from file", apperr.Caller(), err)
		}
		if existingURL.UserID == userID && existingURL.ID == shortURL {
			existingURL.DeletedFlag = true
		}

		urlsToSave = append(urlsToSave, existingURL)
	}

	// Clear file in order to prepare for further encoding
	if err := os.Truncate(r.fileStorage.Name(), 0); err != nil {
		return apperr.NewValueError(fmt.Sprintf("Failed to truncate file: %s", r.fileStorage.Name()), apperr.Caller(), err)
	}

	// Encode urlsToSave to file
	encoder := json.NewEncoder(file)
	for _, url := range urlsToSave {
		err := encoder.Encode(url)
		if err != nil {
			return apperr.NewValueError("unable to encode to file", apperr.Caller(), err)
		}
	}

	return nil
}

// SelectAllByUserID retrieves all URLs by user ID from file
func (r *URLRepository) SelectAllByUserID(ctx context.Context, userID string) ([]model.URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	file, err := os.OpenFile(r.fileStorage.Name(), os.O_RDONLY, perm)
	if err != nil {
		return nil, apperr.NewValueError("unable to open file", apperr.Caller(), err)
	}

	decoder := json.NewDecoder(file)
	defer file.Close()

	urls := make([]model.URL, 0)

	for {
		var url model.URL
		err := decoder.Decode(&url)
		if errors.Is(err, io.EOF) {
			r.logger.Info("Reached end of file while decoding", zap.Error(err))
			break
		}
		if err != nil {
			return nil, apperr.NewValueError("unable to decode from file", apperr.Caller(), err)
		}
		if url.UserID == userID {
			urls = append(urls, url)
		}
	}

	if len(urls) == 0 {
		return nil, apperr.NewValueError(fmt.Sprintf("urls not found by user %s", userID), apperr.Caller(), urlErr.ErrURLNotFound)
	}

	return urls, nil
}

// Insert inserts URL to file
func (r *URLRepository) Insert(ctx context.Context, url model.URL) (*model.URL, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.logger.Info(fmt.Sprintf("Opening file: %s", r.fileStorage.Name()))
	file, err := os.OpenFile(r.fileStorage.Name(), os.O_RDWR|os.O_APPEND, perm)
	if err != nil {
		return nil, apperr.NewValueError("unable to open file", apperr.Caller(), err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(url)
	if err != nil {
		return nil, apperr.NewValueError("unable to encode to file", apperr.Caller(), err)
	}

	return &url, nil
}

// SelectByID retrieves URL from file by ID
func (r *URLRepository) SelectByID(ctx context.Context, key string) (*model.URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	file, err := os.OpenFile(r.fileStorage.Name(), os.O_RDONLY, perm)
	if err != nil {
		return nil, apperr.NewValueError("unable to open file", apperr.Caller(), err)
	}

	decoder := json.NewDecoder(file)
	defer file.Close()

	var url model.URL
	for {
		err := decoder.Decode(&url)
		if errors.Is(err, io.EOF) {
			return nil, apperr.NewValueError(fmt.Sprintf("Url with id %s not found", key), apperr.Caller(), urlErr.ErrURLNotFound)
		}
		if err != nil {
			return nil, apperr.NewValueError("unable to decode from file", apperr.Caller(), err)
		}
		if url.ID == key {
			break
		}
	}

	return &url, nil
}

// SelectAll retrieves all URLs from file
func (r *URLRepository) SelectAll(ctx context.Context) ([]model.URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	file, err := os.OpenFile(r.fileStorage.Name(), os.O_RDONLY, perm)
	if err != nil {
		return nil, apperr.NewValueError("unable to open file", apperr.Caller(), err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	urls := make([]model.URL, 0)

	for {
		var url model.URL
		err := decoder.Decode(&url)
		if errors.Is(err, io.EOF) {
			r.logger.Info("Reached end of file while decoding", zap.Error(err))
			break
		}
		if err != nil {
			return nil, apperr.NewValueError("unable to decode from file", apperr.Caller(), err)
		}
		urls = append(urls, url)
	}

	return urls, nil
}

// DeleteAll deletes all URLs from file
func (r *URLRepository) DeleteAll(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := os.Truncate(r.fileStorage.Name(), 0); err != nil {
		return apperr.NewValueError(fmt.Sprintf("Failed to truncate file: %s", r.fileStorage.Name()), apperr.Caller(), err)
	}
	return nil
}

// Ping pings the file storage
func (r *URLRepository) Ping(ctx context.Context) error {
	file, err := os.OpenFile(r.fileStorage.Name(), os.O_RDONLY, perm)
	if err != nil {
		return apperr.NewValueError("unable to open file", apperr.Caller(), err)
	}
	defer file.Close()
	return nil
}

// InsertAllOrUpdate upserts all URLs to file
func (r *URLRepository) InsertAllOrUpdate(ctx context.Context, urls []model.URL) ([]model.URL, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.logger.Info(fmt.Sprintf("Opening file: %s", r.fileStorage.Name()))
	file, openFileErr := os.OpenFile(r.fileStorage.Name(), os.O_RDWR|os.O_APPEND, perm)
	if openFileErr != nil {
		return nil, apperr.NewValueError("unable to open file", apperr.Caller(), openFileErr)
	}
	defer file.Close()

	// Create map of urls to be saved with id as key
	urlMap := make(map[string]model.URL, len(urls))
	for _, url := range urls {
		urlMap[url.ID] = url
	}

	// Read all urls from file, update with new urls to be saved, store urls to save (with updated ones) in urlsToSave slice
	decoder := json.NewDecoder(file)
	var urlsToSave []model.URL
	for {
		var existingURL model.URL
		err := decoder.Decode(&existingURL)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, apperr.NewValueError("unable to decode from file", apperr.Caller(), err)
		}
		if v, ok := urlMap[existingURL.ID]; ok { // If url exists in map, remove it from map and add it to urlsToSave
			existingURL.Original = v.Original
			existingURL.Shortened = v.Shortened
			delete(urlMap, existingURL.ID)
		}

		urlsToSave = append(urlsToSave, existingURL)
	}

	// Clear file in order to prepare for further encoding
	if err := os.Truncate(r.fileStorage.Name(), 0); err != nil {
		return nil, apperr.NewValueError(fmt.Sprintf("Failed to truncate file: %s", r.fileStorage.Name()), apperr.Caller(), err)
	}

	for _, v := range urlMap {
		urlsToSave = append(urlsToSave, v)
	}

	// Encode urlsToSave to file
	encoder := json.NewEncoder(file)
	for _, url := range urlsToSave {
		err := encoder.Encode(url)
		if err != nil {
			return nil, apperr.NewValueError("unable to encode to file", apperr.Caller(), err)
		}
	}

	// Create map of incomed urls to save with id as key (for faster lookup)
	urlMapIn := make(map[string]model.URL, len(urls))
	for _, url := range urls {
		urlMapIn[url.ID] = url
	}

	// Take from savedURLs slice only urls with id that have been saved
	var savedURLs []model.URL
	for _, url := range urlsToSave {
		if _, ok := urlMapIn[url.ID]; ok {
			savedURLs = append(savedURLs, url)
		}
	}

	return savedURLs, nil
}
