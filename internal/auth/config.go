package core_auth

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

const (
	UserIDContextKey  = "user_id"
	AgentContextKey   = "agent"
	AgentIDContextKey = "agent_id"
)

type Config struct {
	JWTSecret     string        `envconfig:"JWT_SECRET" required:"true"`
	TokenDuration time.Duration `envconfig:"TOKEN_DURATION" required:"true"`
}

func LoadConfig() (*Config, error) {
	var c Config
	if err := envconfig.Process("AUTH", &c); err != nil {
		return nil, fmt.Errorf("load auth config: %w", err)
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
