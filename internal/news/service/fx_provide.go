package service

import (
	"go.uber.org/fx"
)

var FxProvide = fx.Provide(
	NewRealFetcherService,
	NewService,
)
