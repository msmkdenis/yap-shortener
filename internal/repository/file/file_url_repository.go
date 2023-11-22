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

	"github.com/msmkdenis/yap-shortener/internal/apperrors"
	"github.com/msmkdenis/yap-shortener/internal/model"
	"github.com/msmkdenis/yap-shortener/internal/utils"
	"go.uber.org/zap"
)

const (
	perm = 0755
)

type FileURLRepository struct {
	mu          sync.RWMutex
	fileStorage *os.File
	logger      *zap.Logger
}

func NewFileURLRepository(path string, logger *zap.Logger) (*FileURLRepository, error) {
	path = filepath.FromSlash(path)

	dir := filepath.Dir(path)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		logger.Info(fmt.Sprintf("Creating directory: %s", dir))
		err = os.Mkdir(dir, perm)
		if err != nil {
			return nil, apperrors.NewValueError(fmt.Sprintf("Unable to create directory: %s", dir), utils.Caller(), err)
		}
		logger.Info(fmt.Sprintf("Directory %s was created", dir))
	}

	logger.Info(fmt.Sprintf("Creating file: %s", path))
	file, err := os.OpenFile(path, os.O_CREATE, perm)
	if err != nil {
		return nil, apperrors.NewValueError(fmt.Sprintf("Unable to create file: %s", path), utils.Caller(), err)
	}
	logger.Info(fmt.Sprintf("FileStorage %s was created", file.Name()))
	defer file.Close()

	return &FileURLRepository{
		fileStorage: file,
		logger:      logger,
		mu:          sync.RWMutex{},
	}, nil
}

func (r *FileURLRepository) DeleteAllByUserID(ctx context.Context, userID string, shortURLs []string) error {
	return nil
}

func (r *FileURLRepository) SelectAllByUserID(ctx context.Context, userID string) ([]model.URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	file, err := os.OpenFile(r.fileStorage.Name(), os.O_RDONLY, perm)
	if err != nil {
		return nil, apperrors.NewValueError("unable to open file", utils.Caller(), err)
	}

	decoder := json.NewDecoder(file)
	defer file.Close()

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
		if url.UserID == userID {
			urls = append(urls, url)
		}
	}

	if len(urls) == 0 {
		return nil, apperrors.NewValueError(fmt.Sprintf("urls not found by user %s", userID) , utils.Caller(), apperrors.ErrURLNotFound)
	}

	return urls, nil
}

func (r *FileURLRepository) Insert(ctx context.Context, url model.URL) (*model.URL, error) {
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

func (r *FileURLRepository) SelectByID(ctx context.Context, key string) (*model.URL, error) {
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
			return nil, apperrors.NewValueError(fmt.Sprintf("Url with id %s not found", key), utils.Caller(), apperrors.ErrURLNotFound)
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

func (r *FileURLRepository) SelectAll(ctx context.Context) ([]model.URL, error) {
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

func (r *FileURLRepository) DeleteAll(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := os.Truncate(r.fileStorage.Name(), 0); err != nil {
		return apperrors.NewValueError(fmt.Sprintf("Failed to truncate file: %s", r.fileStorage.Name()), utils.Caller(), err)
	}
	return nil
}

func (r *FileURLRepository) Ping(ctx context.Context) error {
	file, err := os.OpenFile(r.fileStorage.Name(), os.O_RDONLY, perm)
	if err != nil {
		return apperrors.NewValueError("unable to open file", utils.Caller(), err)
	}
	defer file.Close()
	return nil
}

func (r *FileURLRepository) InsertAllOrUpdate(ctx context.Context, urls []model.URL) ([]model.URL, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.logger.Info(fmt.Sprintf("Opening file: %s", r.fileStorage.Name()))
	file, err := os.OpenFile(r.fileStorage.Name(), os.O_RDWR|os.O_APPEND, perm)
	if err != nil {
		return nil, apperrors.NewValueError("unable to open file", utils.Caller(), err)
	}
	defer file.Close()

	//Create map of urls to be saved with id as key
	var urlMap = make(map[string]model.URL, len(urls))
	for _, url := range urls {
		urlMap[url.ID] = url
	}

	//Read all urls from file, update with new urls to be saved, store urls to save (with updated ones) in urlsToSave slice
	decoder := json.NewDecoder(file)
	var urlsToSave []model.URL
	for {
		var existingURL model.URL
		err := decoder.Decode(&existingURL)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, apperrors.NewValueError("unable to decode from file", utils.Caller(), err)
		}
		if v, ok := urlMap[existingURL.ID]; ok { //If url exists in map, remove it from map and add it to urlsToSave
			existingURL.Original = v.Original
			existingURL.Shortened = v.Shortened
			delete(urlMap, existingURL.ID)
		}

		urlsToSave = append(urlsToSave, existingURL)
	}

	//Clear file in order to prepare for further encoding
	if err := os.Truncate(r.fileStorage.Name(), 0); err != nil {
		return nil, apperrors.NewValueError(fmt.Sprintf("Failed to truncate file: %s", r.fileStorage.Name()), utils.Caller(), err)
	}

	for _, v := range urlMap {
		urlsToSave = append(urlsToSave, v)
	}

	//Encode urlsToSave to file
	encoder := json.NewEncoder(file)
	for _, url := range urlsToSave {
		err = encoder.Encode(url)
		if err != nil {
			return nil, apperrors.NewValueError("unable to encode to file", utils.Caller(), err)
		}
	}

	//Create map of incomed urls to save with id as key (for faster lookup)
	var urlMapIn = make(map[string]model.URL, len(urls))
	for _, url := range urls {
		urlMapIn[url.ID] = url
	}

	//Take from savedURLs slice only urls with id that have been saved
	var savedURLs []model.URL
	for _, url := range urlsToSave {
		if _, ok := urlMapIn[url.ID]; ok {
			savedURLs = append(savedURLs, url)
		}
	}

	return savedURLs, nil
}
