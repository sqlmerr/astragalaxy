package service

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Seed int64 `required:"true"`
}

func NewConfig() (Config, error) {
	var cfg Config

	if err := envconfig.Process("GAME", &cfg); err != nil {
		return Config{}, fmt.Errorf("load game config: %w", err)
	}

	return cfg, nil
}

func NewConfigMust() Config {
	cfg, err := NewConfig()
	if err != nil {
		panic(err)
	}
	return cfg
}
