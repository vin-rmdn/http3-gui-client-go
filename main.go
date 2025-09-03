package main

import (
	"crypto/tls"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/vin-rmdn/http3-gui-client-go/internal/application"
	"github.com/vin-rmdn/http3-gui-client-go/internal/config"
	"github.com/vin-rmdn/http3-gui-client-go/internal/logger"
)

func main() {
	conf, err := config.Load("./config.yml")
	if err != nil {
		slog.Error("cannot setup configuration", slog.String("error", err.Error()))

		os.Exit(1) // TODO: get a better code
	}

	log, err := logger.Setup(*conf)
	if err != nil {
		slog.Error("cannot unmarshal log level", slog.String("level", conf.Log.Level))

		const exitCodeDataError = 65
		os.Exit(exitCodeDataError)
	}

	log.Info("Hello, world!", slog.String("todays_message", conf.Today.Message))

	const title = "HTTP3 GUI Client in Go"
	glib.SetApplicationName(title) // Linux
	glib.SetPrgname(title)         // Darwin

	gtk.Init(nil)

	gtkApp, err := gtk.ApplicationNew("dev.systrshr.http3_client", glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		slog.Error("cannot initialize application", slog.String("error", err.Error()))
		os.Exit(1)
	}

	const signalActivate = "activate"
	gtkApp.Connect(signalActivate, func() {
		app := &application.Application{
			Application: gtkApp,
			HTTPClient: &http.Client{
				Transport: &http3.Transport{
					TLSClientConfig: &tls.Config{
						NextProtos: []string{http3.NextProtoH3},
					},
					QUICConfig: &quic.Config{},
				},
				Timeout: time.Duration(conf.HTTP.TimeoutInSecond) * time.Second,
			},
		}

		// TODO: send an exit signal and tidy up if error is detected
		app.Activate(conf)
	})

	exitCode := gtkApp.Run(os.Args)

	os.Exit(exitCode)
}
