package storage

import (
	"errors"
	"sync"

	"github.com/aidlatyp/ya-pr-shortener/internal/app/domain"
)

type URLMemoryStorage struct {
	storage map[string]string
	mutex   sync.RWMutex
}

func NewURLMemoryStorage() *URLMemoryStorage {
	return &URLMemoryStorage{
		storage: make(map[string]string),
	}
}

func (u *URLMemoryStorage) Store(url *domain.URL) error {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	u.storage[(*url).Short] = (*url).Orig
	return nil
}

func (u *URLMemoryStorage) FindByKey(key string) (*domain.URL, error) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()
	orig, ok := u.storage[key]
	if !ok {
		return nil, errors.New("key does not exist")
	}
	url := domain.NewURL(orig, key)
	return url, nil
}
