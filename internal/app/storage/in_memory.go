package storage

import (
	"errors"
	"sync"

	"github.com/aidlatyp/ya-pr-shortener/internal/app/domain"
)

type URLMemoryStorage struct {
	storage   map[string]string
	mutex     sync.RWMutex
	userLinks map[string][]string
}

func NewURLMemoryStorage() *URLMemoryStorage {
	return &URLMemoryStorage{
		storage:   make(map[string]string),
		userLinks: make(map[string][]string),
	}
}

func (u *URLMemoryStorage) Store(url *domain.URL) error {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	u.storage[url.Short] = url.Orig

	if url.Owner != "" {
		u.userLinks[url.Owner] = append(u.userLinks[url.Owner], url.Short)
	}
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

func (u *URLMemoryStorage) FindAll(userKey string) []*domain.URL {
	u.mutex.RLock()
	defer u.mutex.RUnlock()
	userBucket, ok := u.userLinks[userKey]
	if !ok {
		return nil
	}

	resultList := make([]*domain.URL, 0, len(userBucket))

	for _, key := range userBucket {
		orig, _ := u.storage[key]
		url := domain.NewURL(orig, key)

		resultList = append(resultList, url)
	}
	return resultList
}
