package pkg

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Steps []Step `yaml:"steps"`
}

type Step struct {
	Name     string                 `yaml:"name"`
	Template string                 `yaml:"template"`
	Option   map[string]interface{} `yaml:"option"`
}

func NewConfig(path string) (Config, error) {
	b, e := os.ReadFile(path)
	if e != nil {
		return Config{}, e
	}

	var config Config
	if e := yaml.Unmarshal(b, &config); e != nil {
		return Config{}, e
	}
	return config, nil
}
