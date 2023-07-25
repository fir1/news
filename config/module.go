package config

import "go.uber.org/fx"

var FxProvide = fx.Provide(
	NewParsedConfig,
)
