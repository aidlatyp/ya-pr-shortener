package storage

type MemoryStorage struct {
	shortURLs map[string]string
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		shortURLs: map[string]string{
			//"https://yatest.ru": "test", // need for tests
		},
	}
}

func (ms *MemoryStorage) GetAll() (map[string]string, error) {
	return ms.shortURLs, nil
}

func (ms *MemoryStorage) GetByKey(key string) (string, error) {
	return ms.shortURLs[key], nil
}

func (ms *MemoryStorage) Set(key string, value string) error {
	ms.shortURLs[key] = value
	return nil
}

func (ms *MemoryStorage) CloseResources() error {
	return nil
}
