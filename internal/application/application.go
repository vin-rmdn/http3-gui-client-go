package application

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/gotk3/gotk3/gtk"
	"github.com/quic-go/quic-go/http3"
	"github.com/vin-rmdn/http3-gui-client-go/internal/config"
)

type Application struct {
	*gtk.Application

	URLTextView        *gtk.TextView
	MethodComboBoxText *gtk.ComboBoxText
	HTTPClient         *http.Client
}

func (a *Application) Activate(conf *config.Configuration) {
	window, err := gtk.ApplicationWindowNew(a.Application)
	if err != nil {
		slog.Error("cannot initialize new application window", slog.String("error", err.Error()))
		os.Exit(1)
	}

	window.SetTitle(conf.Window.Title)
	window.SetDefaultSize(800, 600)

	window.Connect("destroy", func() {
		a.HTTPClient.Transport.(*http3.Transport).Close()
		gtk.MainQuit()
	})

	verticalGrid, _ := gtk.GridNew()
	verticalGrid.SetOrientation(gtk.ORIENTATION_VERTICAL)
	verticalGrid.SetBorderWidth(10)
	verticalGrid.SetRowSpacing(10)

	verticalGrid.Add(a.createHTTPURLInput())

	sendRequestButton, _ := gtk.ButtonNewWithLabel("Send request")
	sendRequestButton.SetHExpand(false)

	const signalClicked = "clicked"
	sendRequestButton.Connect(signalClicked, a.sendHTTP())
	verticalGrid.Add(sendRequestButton)

	window.Add(verticalGrid)

	window.ShowAll()
}

func (a *Application) createHTTPURLInput() gtk.IWidget {
	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)

	a.MethodComboBoxText, _ = gtk.ComboBoxTextNew()
	a.MethodComboBoxText.SetHExpand(false)

	// TODO: provide all methods
	supportedHTTPMethods := []string{
		http.MethodGet,
		http.MethodPost,
	}
	for _, method := range supportedHTTPMethods {
		a.MethodComboBoxText.AppendText(method)
	}
	a.MethodComboBoxText.SetActive(0)

	box.PackStart(a.MethodComboBoxText, false, false, 5)

	a.URLTextView, _ = gtk.TextViewNew()
	a.URLTextView.SetHExpand(true)
	a.URLTextView.SetTopMargin(10)
	textBoxBuffer, _ := a.URLTextView.GetBuffer()
	textBoxBuffer.SetText("Enter URL here")

	box.PackStart(a.URLTextView, true, true, 5)

	return box
}

// sendHTTP gets the appropriate value from the boxes and creates a HTTP connection.
// TODO: extract this to its own package
func (a *Application) sendHTTP() func() {
	return func() {
		urlBuffer, _ := a.URLTextView.GetBuffer()
		start, end := urlBuffer.GetBounds()

		url, _ := urlBuffer.GetText(start, end, false)

		method := a.MethodComboBoxText.GetActiveText()

		slog.Info(
			"ready to trigger http request",
			slog.String("url", url),
			slog.String("method", method),
		)

		// TODO: support request body
		req, err := http.NewRequestWithContext(context.Background(), method, url, http.NoBody)
		if err != nil {
			slog.Error("cannot create request", slog.String("error", err.Error()))
			return
		}

		response, err := a.HTTPClient.Do(req)
		if err != nil {
			slog.Error("cannot execute http call", slog.String("error", err.Error()))
			return
		}

		// response.Status
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			slog.Error("cannot decode response", slog.String("error", err.Error()))
			return
		}

		defer func() {
			_ = response.Body.Close()
		}()

		slog.Info(
			"response received",
			slog.String("status", response.Status),
			slog.String("body", string(responseBody)),
			slog.String("proto_version", response.Proto),
		)
	}
}
