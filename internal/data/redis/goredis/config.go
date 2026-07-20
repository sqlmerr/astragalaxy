package redis_goredis

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Username string
	Password string
	Addr     string
	DB       int
}

func LoadConfig() (*Config, error) {
	var c Config
	if err := envconfig.Process("REDIS", &c); err != nil {
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
