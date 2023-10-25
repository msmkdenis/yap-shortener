package file

import (
	"encoding/json"
	"github.com/msmkdenis/yap-shortener/internal/model"
	"go.uber.org/zap"
	"io"
	"log"
	"os"
	"path/filepath"
)

type FileURLRepository struct {
	fileStoragePath string
	logger          *zap.Logger
}

func NewFileURLRepository(path string, logger *zap.Logger) *FileURLRepository {
	path = filepath.FromSlash(path)

	dir := filepath.Dir(path)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			log.Fatalf("Error: %s", err)
			return nil
		}
	}

	return &FileURLRepository{
		fileStoragePath: path,
		logger:          logger,
	}
}

func (r *FileURLRepository) Insert(url model.URL) (*model.URL, error) {
	if savedURL, err := r.SelectByID(url.ID); err == nil {
		return savedURL, nil
	}

	r.logger.Info("Opening file")
	file, err := os.OpenFile(r.fileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalf("Ошибка при открытии %s", err)
		return nil, err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)

	err = encoder.Encode(url)

	if err != nil {
		r.logger.Info("Could not encode url to file")
		return nil, err
	}

	return &url, nil
}

func (r *FileURLRepository) SelectByID(key string) (*model.URL, error) {

	file, err := os.OpenFile(r.fileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalf("Ошибка при открытии %s", err)
		return nil, err
	}

	decoder := json.NewDecoder(file)

	defer file.Close()

	var url model.URL

	for {
		err := decoder.Decode(&url)
		if err == io.EOF {
			r.logger.Info("Could not find url in file")
			return nil, err
		}
		if url.ID == key {
			break
		}
	}

	return &url, nil
}

func (r *FileURLRepository) SelectAll() ([]model.URL, error) {
	file, err := os.OpenFile(r.fileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalf("Ошибка при открытии %s", err)
		return nil, err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	var answer []model.URL

	for {
		var url model.URL
		err := decoder.Decode(&url)
		if err == io.EOF {
			r.logger.Info("Reached end of file while decoding")
			break
		}
		answer = append(answer, url)
	}

	return answer, nil
}

func (r *FileURLRepository) DeleteAll() {
	if err := os.Truncate(r.fileStoragePath, 0); err != nil {
		log.Printf("Failed to truncate: %v", err)
	}
}
