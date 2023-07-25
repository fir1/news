package cache

// CacheClientInterface represents a allegro/bigcache client
type CacheClientInterface interface {
	Get(key string) ([]byte, error)
	Set(key string, entry []byte) error
	Delete(key string) error
	Reset() error
}
