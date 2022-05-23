package storage

import "sync"

type MemoryStorage struct {
	mutex     sync.RWMutex
	shortURLs map[string]string
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		shortURLs: make(map[string]string),
	}
}

func (ms *MemoryStorage) GetAll() (map[string]string, error) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	return ms.shortURLs, nil
}

func (ms *MemoryStorage) GetByKey(key string) (string, error) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	return ms.shortURLs[key], nil
}

func (ms *MemoryStorage) Set(key string, value string) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	ms.shortURLs[key] = value
	return nil
}

func (ms *MemoryStorage) CloseResources() error {
	return nil
}
