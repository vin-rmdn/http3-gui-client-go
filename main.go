package main

import (
	"log/slog"
	"os"

	"github.com/gotk3/gotk3/gtk"
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

	gtk.Init(nil)

	window, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		slog.Error("cannot create new window", slog.String("error", err.Error()))
	}

	window.SetTitle(conf.Window.Title)
	window.Connect("destroy", func() {
		gtk.MainQuit()
	})

	topLabel, err := gtk.LabelNew("This is a label!")
	if err != nil {
		slog.Error("cannot create a new label")
	}

	window.Add(topLabel)

	window.SetDefaultSize(800, 600)
	window.ShowAll()

	gtk.Main()
}
