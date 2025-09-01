package main

import (
	"log/slog"
	"os"

	"github.com/vin-rmdn/http3-gui-client-go/internal/config"
	"github.com/vin-rmdn/http3-gui-client-go/internal/logger"
)

func main() {
	conf := config.Load("./config.yml")

	log, err := logger.Setup(conf)
	if err != nil {
		slog.Error(
			"cannot unmarshal log level",
			slog.String("level", conf.Log.Level),
		)

		const exitCodeDataError = 65
		os.Exit(exitCodeDataError)
	}
	log.Info("Hello, world!", slog.String("todays_message", conf.Today.Message))
}
