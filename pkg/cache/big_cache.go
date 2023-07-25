package cache

import (
	"context"
	"github.com/allegro/bigcache/v3"
	"time"
)

type Bigcache struct {
	client *bigcache.BigCache
}

func NewBigcache() (CacheClientInterface, error) {
	clientCache, err := bigcache.New(context.Background(), bigcache.DefaultConfig(5*time.Minute))
	if err != nil {
		return nil, err
	}

	return Bigcache{
		client: clientCache,
	}, nil
}

func (b Bigcache) Get(key string) ([]byte, error) {
	return b.client.Get(key)
}

func (b Bigcache) Set(key string, entry []byte) error {
	return b.client.Set(key, entry)
}

func (b Bigcache) Delete(key string) error {
	return b.client.Delete(key)
}

func (b Bigcache) Reset() error {
	return b.client.Reset()
}
