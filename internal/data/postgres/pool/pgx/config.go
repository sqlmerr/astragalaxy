package pgx_pool

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	User     string
	Password string
	Host     string
	Port     int           `default:"5432"`
	Database string        `envconfig:"db"`
	Timeout  time.Duration `default:"5s"`
}

func LoadConfig() (*Config, error) {
	var c Config
	if err := envconfig.Process("POSTGRES", &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func LoadConfigMust() *Config {
	c, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	return c
}
