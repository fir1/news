package http

import (
	"github.com/fir1/news/config"
	newsSvc "github.com/fir1/news/internal/news/service"
	"github.com/fir1/news/pkg/cache"
	"github.com/go-chi/chi/v5"

	"github.com/sirupsen/logrus"
)

type Service struct {
	router            *chi.Mux
	logger            *logrus.Logger
	stockSymbol       string
	stockNumberOfDays int
	newsService       newsSvc.NewsInterface
	config            config.Config
	cacheClient       cache.CacheClientInterface
}

func NewService(logger *logrus.Logger,
	newsSvc newsSvc.NewsInterface,
	cnf config.Config,
	cc cache.CacheClientInterface,
) *Service {
	return &Service{
		logger:      logger,
		newsService: newsSvc,
		config:      cnf,
		cacheClient: cc,
	}
}
