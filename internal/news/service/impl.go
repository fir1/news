package service

import "github.com/fir1/news/pkg/cache"

type Service struct {
	NewsFetcher NewsFetcher
	cacheClient cache.CacheClientInterface
}

func NewService(nf NewsFetcher,
	cc cache.CacheClientInterface,
) NewsInterface {
	return Service{
		NewsFetcher: nf,
		cacheClient: cc,
	}
}
