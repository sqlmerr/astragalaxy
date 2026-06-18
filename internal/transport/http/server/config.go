package http_server

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Addr            string        `required:"true"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" required:"true"`
}

func LoadConfig() (*Config, error) {
	var c Config
	if err := envconfig.Process("HTTP", &c); err != nil {
		return nil, fmt.Errorf("load http config: %w", err)
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
