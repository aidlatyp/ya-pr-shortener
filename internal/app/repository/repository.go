package repository

import (
	"errors"
	"github.com/aidlatyp/ya-pr-shortener/internal/app/domain"
	"sync"
)

type URLRepo struct {
	storage map[string]string
	mutex   sync.RWMutex
}

func NewURLRepo() *URLRepo {
	return &URLRepo{
		storage: make(map[string]string),
	}
}

func (u *URLRepo) Store(url *domain.URL) error {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	u.storage[(*url).Short] = (*url).Orig
	return nil
}

func (u *URLRepo) FindByKey(key string) (*domain.URL, error) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()
	orig, ok := u.storage[key]
	if !ok {
		return nil, errors.New("key does not exist")
	}
	url := domain.NewURL(orig, key)
	return url, nil
}
