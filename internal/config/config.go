package config

import (
	"log/slog"
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

	Window struct {
		Title string
	}
}

func Load(file string) (*Configuration, error) {
	buffer, err := os.ReadFile(file)
	if err != nil {
		slog.Error("cannot read config.yml file")

		return nil, err
	}

	config := &Configuration{}
	if err := yaml.Unmarshal(buffer, config); err != nil {
		slog.Error("cannot unmarshal config file")

		return nil, err
	}

	return config, nil
}
