package repositories

type Repository interface {
	GetAll() (map[string]string, error)
	GetByKey(key string) (string, error)
	Set(key string, value string) error
	CloseResources() error
}
