package main

import (
	"crypto/tls"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/vin-rmdn/http3-gui-client-go/internal/view"
	"github.com/vin-rmdn/http3-gui-client-go/internal/config"
	"github.com/vin-rmdn/http3-gui-client-go/internal/controller"
	"github.com/vin-rmdn/http3-gui-client-go/internal/logger"
)

func main() {
	// Patch for macOS bundling
	execPath, err := os.Executable()
	if err != nil {
		slog.Error("cannot get executable path", slog.String("error", err.Error()))
	}

	execPath, _ = filepath.EvalSymlinks(execPath)

	conf, err := config.Load(
		"./config.yml",
		filepath.Join(filepath.Dir(execPath), "..", "Resources", "config.yml"),
	)
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

	gtk.Init()

	gtkApp := gtk.NewApplication("dev.systrshr.http3_client", gio.ApplicationFlagsNone)
	if err != nil {
		slog.Error("cannot initialize application", slog.String("error", err.Error()))
		os.Exit(1)
	}

	app := &view.View{
		Application: gtkApp,
	}
	httpClient := &http.Client{
		Transport: &http3.Transport{
			TLSClientConfig: &tls.Config{
				NextProtos: []string{http3.NextProtoH3},
			},
			QUICConfig: &quic.Config{},
		},
		Timeout: time.Duration(conf.HTTP.TimeoutInSecond) * time.Second,
	}

	const signalActivate = "activate"
	gtkApp.Connect(signalActivate, func() {
		// TODO: send an exit signal and tidy up if error is detected
		if err := app.Activate(conf); err != nil {
			slog.Error("cannot activate application", slog.String("error", err.Error()))
			os.Exit(1)
		}

		c := controller.New(httpClient, app)
		_ = c
	})

	exitCode := gtkApp.Run(os.Args)

	os.Exit(exitCode)
}
