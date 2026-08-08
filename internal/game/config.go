package game

import (
	"fmt"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var (
	k = koanf.New(".")
)

type RulesConfig struct {
	DisableCooldowns      bool `koanf:"disableCooldowns"`
	DisableInventoryLimit bool `koanf:"disableInventoryLimit"`
}

type Config struct {
	Seed  int64       `koanf:"seed"`
	Rules RulesConfig `koanf:"rules"`
}

func LoadConfig() (Config, error) {
	if err := k.Load(file.Provider("config.yaml"), yaml.Parser()); err != nil {
		return Config{}, fmt.Errorf("load yaml config: %w", err)
	}

	var config Config
	if err := k.Unmarshal("", &config); err != nil {
		return Config{}, fmt.Errorf("unmarshal config: %w", err)
	}

	return config, nil
}

func LoadConfigMust() Config {
	cfg, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	return cfg
}
