package application

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gotk3/gotk3/gtk"
	"github.com/vin-rmdn/http3-gui-client-go/internal/config"
)

const signalDestroy = "destroy"

type View struct {
	*gtk.Application

	Window             *gtk.ApplicationWindow
	URLTextView        *gtk.TextView
	MethodComboBoxText *gtk.ComboBoxText
	SendRequestButton  *gtk.Button
}

func (a *View) Activate(conf *config.Configuration) error {
	window, err := gtk.ApplicationWindowNew(a.Application)
	if err != nil {
		return fmt.Errorf("cannot initialize new application window: %w", err)
	}

	a.Window = window
	a.Window.SetTitle(conf.Window.Title)
	a.Window.SetDefaultSize(800, 600)

	a.Window.ConnectAfter(signalDestroy, func() {
		a.Quit()
	})

	verticalGrid, _ := gtk.GridNew()
	verticalGrid.SetOrientation(gtk.ORIENTATION_VERTICAL)
	verticalGrid.SetBorderWidth(10)
	verticalGrid.SetRowSpacing(10)

	verticalGrid.Add(a.createHTTPURLInput())

	a.SendRequestButton, err = gtk.ButtonNewWithLabel("Send request")
	if err != nil {
		return fmt.Errorf("cannot create send request button: %w", err)
	}

	a.SendRequestButton.SetHExpand(false)

	verticalGrid.Add(a.SendRequestButton)

	window.Add(verticalGrid)

	window.ShowAll()

	return nil
}

func (a *View) SetOnSendRequestFunction(callback func(*http.Request)) {
	const signalClicked = "clicked"
	a.SendRequestButton.Connect(signalClicked, func() {
		urlBuffer, _ := a.URLTextView.GetBuffer()
		start, end := urlBuffer.GetBounds()

		url, _ := urlBuffer.GetText(start, end, false)

		method := a.MethodComboBoxText.GetActiveText()

		slog.Debug(
			"ready to trigger http request",
			slog.String("url", url),
			slog.String("method", method),
		)

		// TODO: support request body
		request, err := http.NewRequestWithContext(context.Background(), method, url, http.NoBody)
		if err != nil {
			slog.Error("cannot create request", slog.String("error", err.Error()))
			return
		}

		callback(request)
	})
}

func (a *View) createHTTPURLInput() gtk.IWidget {
	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)

	a.MethodComboBoxText, _ = gtk.ComboBoxTextNew()
	a.MethodComboBoxText.SetHExpand(false)

	supportedHTTPMethods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodConnect,
		http.MethodHead,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodOptions,
		http.MethodTrace,
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

func (a *View) SetDestroyFunction(callback func()) {
	a.Window.Connect(signalDestroy, callback)
}
