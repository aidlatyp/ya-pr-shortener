package storage

import (
	"errors"
	"sync"

	"github.com/aidlatyp/ya-pr-shortener/internal/app/domain"
)

type URLMemoryStorage struct {
	linksStorage map[uniqID]string
	mutex        sync.RWMutex
	userLinks    map[string][]uniqID
}

type uniqID string

func newURLMemoryStorage() *URLMemoryStorage {
	// Do not have duplications of URLs, full maps scan for any use cases, etc,
	// userLinks is map of slices, each slice is a list of refs (keys)
	// to a user links in linksStorage map. To make this ref principle more explicit,
	// redundant [uniqID] type declared
	return &URLMemoryStorage{
		userLinks:    make(map[string][]uniqID),
		linksStorage: make(map[uniqID]string),
	}
}

func (u *URLMemoryStorage) BatchWrite(urls []domain.URL) error {
	for _, v := range urls {
		// reserved error api
		_ = u.Store(&v)
	}
	return nil
}

func (u *URLMemoryStorage) Store(url *domain.URL) error {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	u.linksStorage[uniqID(url.Short)] = url.Orig

	if url.Owner != "" {
		u.userLinks[url.Owner] = append(u.userLinks[url.Owner], uniqID(url.Short))
	}
	return nil
}

func (u *URLMemoryStorage) FindByKey(key string) (*domain.URL, error) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	orig, ok := u.linksStorage[uniqID(key)]
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
		orig := u.linksStorage[key]
		url := domain.NewURL(orig, string(key))
		resultList = append(resultList, url)
	}

	return resultList
}
