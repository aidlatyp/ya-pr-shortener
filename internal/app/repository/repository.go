package repository

import (
	"errors"
	"github.com/aidlatyp/ya-pr-shortener/internal/app/domain"
	"sync"
)

type UrlRepo struct {
	storage map[string]string
	mutex   sync.RWMutex
}

func NewUrlRepo() *UrlRepo {
	return &UrlRepo{
		storage: make(map[string]string),
	}
}

func (u *UrlRepo) Store(url *domain.Url) error {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	u.storage[(*url).Short] = (*url).Orig
	return nil
}

func (u *UrlRepo) FindByKey(key string) (*domain.Url, error) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()
	orig, ok := u.storage[key]
	if !ok {
		return nil, errors.New("key does not exist")
	}
	url := domain.NewUrl(orig, key)
	return url, nil
}
