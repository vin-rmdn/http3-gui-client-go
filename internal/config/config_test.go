package config_test

import (
	"testing"

	"github.com/vin-rmdn/http3-gui-client-go/internal/config"
)

func TestConfig(t *testing.T) {
	emptyConfig := &config.Configuration{}

	conf, err := config.Load("../../config.sample.yml")
	if err != nil {
		t.Errorf("expected no error, got one: %s", err.Error())
	}

	if conf == emptyConfig {
		t.Error("expected populated config, got nothing")
	}
}
