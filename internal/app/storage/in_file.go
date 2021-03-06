package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aidlatyp/ya-pr-shortener/internal/app/domain"
	"github.com/aidlatyp/ya-pr-shortener/internal/app/usecase"
)

const LineBreak byte = '\n'

type PersistentStorage struct {
	cache usecase.Repository
	file  *os.File
}

func NewStorage(path string) usecase.Repository {

	var store usecase.Repository = newURLMemoryStorage()
	if path != "" {
		persistentStorage, err := newPersistentStorage(path, store)
		if err != nil {
			log.Fatalf("filepath set, but can't start in persistent mode %v ", err.Error())
		}
		store = persistentStorage
	}
	return store
}

func newPersistentStorage(path string, cache usecase.Repository) (*PersistentStorage, error) {

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	sc := bufio.NewScanner(file)
	for sc.Scan() {

		var url *domain.URL
		line := sc.Bytes()
		err = json.Unmarshal(line, &url)
		if err != nil {
			return nil, fmt.Errorf("error while unmarshal from file %v ", err)
		}

		err = cache.Store(url)
		if err != nil {
			return nil, fmt.Errorf("error while filling cache %v ", err)
		}
	}

	return &PersistentStorage{
		file:  file,
		cache: cache,
	}, nil
}

func (p *PersistentStorage) Store(url *domain.URL) error {
	bytes, err := json.Marshal(url)
	if err != nil {
		return fmt.Errorf("error while marshaling data  %v ", err)
	}
	bytes = append(bytes, LineBreak)
	_, err = p.file.Write(bytes)
	if err != nil {
		return fmt.Errorf("error while writing to file %v ", err)
	}
	err = p.cache.Store(url)
	return err
}

func (p *PersistentStorage) FindByKey(key string) (*domain.URL, error) {
	return p.cache.FindByKey(key)
}

func (p *PersistentStorage) FindAll(key string) []*domain.URL {
	return p.cache.FindAll(key)
}

func (p *PersistentStorage) BatchWrite(urls []domain.URL) error {
	return p.cache.BatchWrite(urls)
}

func (p *PersistentStorage) Close() error {
	return p.file.Close()
}
