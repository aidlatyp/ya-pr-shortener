package storage

import (
	"errors"
	"github.com/aidlatyp/ya-pr-shortener/internal/app/domain"
	"sync"
)

type URLStorage struct {
	storage map[string]string
	mutex   sync.RWMutex
}

func NewURLStorage() *URLStorage {
	return &URLStorage{
		storage: make(map[string]string),
	}
}

func (u *URLStorage) Store(url *domain.URL) error {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	u.storage[(*url).Short] = (*url).Orig
	return nil
}

func (u *URLStorage) FindByKey(key string) (*domain.URL, error) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()
	orig, ok := u.storage[key]
	if !ok {
		return nil, errors.New("key does not exist")
	}
	url := domain.NewURL(orig, key)
	return url, nil
}
