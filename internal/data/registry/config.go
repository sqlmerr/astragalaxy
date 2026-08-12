package registry

import (
	"fmt"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var k = koanf.New(".")

type Config struct {
	ItemsPath      string `koanf:"itemsPath"`
	ResourcesPath  string `koanf:"resourcesPath"`
	RecipesPath    string `koanf:"recipesPath"`
	FacilitiesPath string `koanf:"facilitiesPath"`
}

func LoadConfig() (Config, error) {
	if err := k.Load(file.Provider("config.yaml"), yaml.Parser()); err != nil {
		return Config{}, fmt.Errorf("load yaml config: %w", err)
	}

	var config Config
	if err := k.Unmarshal("registry", &config); err != nil {
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
