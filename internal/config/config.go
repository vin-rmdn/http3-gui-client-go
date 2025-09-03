package config

import (
	"errors"
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

	HTTP struct {
		TimeoutInSecond int
	}
}

func Load(files ...string) (*Configuration, error) {
	var errs error
	var succeeds bool

	config := &Configuration{}

	for _, file := range files {
		buffer, err := os.ReadFile(file)
		if err != nil {
			slog.Error("cannot read config.yml file")

			errs = errors.Join(errs, err)
			continue
		}

		if err := yaml.Unmarshal(buffer, config); err != nil {
			slog.Error("cannot unmarshal config file")

			errs = errors.Join(errs, err)
			continue
		}

		succeeds = true

		break
	}

	if !succeeds {
		return nil, errs
	}

	return config, nil
}
