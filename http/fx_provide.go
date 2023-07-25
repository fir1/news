package http

import (
	"github.com/fir1/news/pkg/cache"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"os"
)

var FxProvide = fx.Provide(
	NewTextLogger,
	NewService,
	cache.NewBigcache,
)

func NewTextLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		TimestampFormat: "2006-01-02 15:04:05.999999999",
		FullTimestamp:   true,
	})
	return logger
}
