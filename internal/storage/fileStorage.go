package storage

import (
	"github.com/aidlatyp/ya-pr-shortener/internal/service"
)

type FileStorage struct {
	writer    *service.Writer
	reader    *service.Reader
	shortURLs map[string]string
}

func NewFileStorage(filePath string) (*FileStorage, error) {

	w, err := service.NewWriter(filePath)
	if err != nil {
		return nil, err
	}

	r, err := service.NewReader(filePath)
	if err != nil {
		return nil, err
	}

	// читаем один раз, потом работаем в памяти
	shortURLs, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	return &FileStorage{
		writer:    w,
		reader:    r,
		shortURLs: shortURLs,
	}, nil
}

func (f *FileStorage) CloseResources() error {
	return f.writer.Close()
}

func (f *FileStorage) GetAll() (map[string]string, error) {
	return f.shortURLs, nil
}

func (f *FileStorage) GetByKey(key string) (string, error) {
	records, err := f.GetAll()
	if err != nil {
		return "", err
	}
	return records[key], nil
}

func (f *FileStorage) Set(key string, value string) error {

	f.shortURLs[key] = value

	r := &service.Record{Key: key, Value: value}
	if err := f.writer.Write(r); err != nil {
		return err
	}
	return nil
}
