package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func NewParsedConfig() (Config, error) {
	_ = godotenv.Load() // The Original .env
	cnf := Config{}
	err := envconfig.Process("", &cnf)
	return cnf, err
}

type Config struct {
	Environment          string `envconfig:"ENVIRONMENT" default:"dev"`
	ServerHostName       string `envconfig:"SERVER_HOST_NAME" default:"http://0.0.0.0"`
	Port                 int    `envconfig:"PORT" default:"8080"`
	LoadBalancerHostPort int    `envconfig:"LOAD_BALANCER_HOST_PORT" default:"8080"`
}
