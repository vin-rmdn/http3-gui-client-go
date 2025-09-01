package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Today struct {
		Message string
	}

	Log struct {
		Level string
	}
}

func Load(file string) Configuration {
	buffer, err := os.ReadFile(file)
	if err != nil {
		panic("cannot read config.yml file")
	}

	var config Configuration
	if err := yaml.Unmarshal(buffer, &config); err != nil {
		panic("cannot unmarshal config file")
	}

	return config
}
