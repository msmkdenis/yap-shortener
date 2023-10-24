package file

import (
	"encoding/json"
	"fmt"
	"github.com/msmkdenis/yap-shortener/internal/model"
	"io"
	"log"
	"os"
	"path/filepath"
)

type FileURLRepository struct {
	fileStoragePath string
}

func NewFileURLRepository(path string) *FileURLRepository {
	/*	dir := filepath.Dir(path)

		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err = os.Mkdir(dir, 0755)
			if err != nil {
				log.Fatalf("Error: %s", err)
				return nil
			}
		}*/
	dir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	} else {
		log.Println(dir)
	}
	path = filepath.FromSlash(path)
	fmt.Println("++++++++++++++++++++++++")
	fmt.Println(path)
	return &FileURLRepository{fileStoragePath: dir + path}
}

func (r *FileURLRepository) Insert(url model.URL) (*model.URL, error) {
	fmt.Println(r.fileStoragePath)

	file, err := os.OpenFile(r.fileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalf("Ошибка при открытии %s", err)
		return nil, err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)

	err = encoder.Encode(url)

	if err != nil {
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
			return nil, err
		}
		if err != nil {
			log.Fatalf("Error decoding: %v", err)
		}

		if url.ID == key {
			break
		}
	}

	return &url, nil
}

func (r *FileURLRepository) SelectAll() ([]string, error) {
	file, err := os.OpenFile(r.fileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalf("Ошибка при открытии %s", err)
		return nil, err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	var answer []string

	for {
		var url model.URL
		err := decoder.Decode(&url)

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Error decoding: %v", err)
		}
		answer = append(answer, url.ID+" "+url.Original+" "+url.Shortened+"\n")
	}
	return answer, nil
}

func (r *FileURLRepository) DeleteAll() {
	if err := os.Truncate(r.fileStoragePath, 0); err != nil {
		log.Printf("Failed to truncate: %v", err)
	}
}
