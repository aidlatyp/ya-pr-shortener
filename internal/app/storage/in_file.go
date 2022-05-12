package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aidlatyp/ya-pr-shortener/internal/app/domain"
)

const LineBreak byte = '\n'

type PersistentStorage struct {
	cache *URLMemoryStorage
	file  *os.File
}

func NewPersistentStorage(path string) (*PersistentStorage, error) {

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}

	cache := NewURLMemoryStorage()
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
		return err
	}

	bytes = append(bytes, LineBreak)
	_, err = p.file.Write(bytes)
	if err != nil {
		return err
	}

	p.cache.Store(url)
	return nil
}

func (p *PersistentStorage) FindByKey(key string) (*domain.URL, error) {
	url, err := p.cache.FindByKey(key)
	if err != nil {
		return nil, err
	}
	return url, nil
}

func (p *PersistentStorage) Close() error {
	return p.file.Close()
}
