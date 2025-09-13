package view

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"strings"

	gtk "github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/vin-rmdn/http3-gui-client-go/internal/config"
)

const signalDestroy = "destroy"

type View struct {
	*gtk.Application

	Logger *slog.Logger

	Window              *gtk.ApplicationWindow
	URLTextView         *gtk.TextView
	MethodDropDown      *gtk.DropDown
	MethodList          []string
	RequestBodyTextView *gtk.TextView
	SendRequestButton   *gtk.Button
	ResponseTextView    *gtk.TextView
}

func (a *View) Activate(conf *config.Configuration) error {
	a.Window = gtk.NewApplicationWindow(a.Application)
	a.Window.SetTitle(conf.Window.Title)
	a.Window.SetDefaultSize(800, 600)

	a.Window.ConnectAfter(signalDestroy, func() {
		a.Application.Quit()
	})

	a.RequestBodyTextView = gtk.NewTextView()
	a.RequestBodyTextView.SetHExpand(true)
	a.RequestBodyTextView.SetSizeRequest(-1, 100)

	a.SendRequestButton = gtk.NewButtonWithLabel("Send request")
	a.SendRequestButton.SetHExpand(false)

	httpUIGrid := a.addAllToVerticalGrid(
		a.createHTTPURLInput(),
		func() gtk.Widgetter {
			requestBodyLabel := gtk.NewLabel("Request body:")
			requestBodyLabel.SetHAlign(gtk.AlignStart)

			return requestBodyLabel
		}(),
		a.RequestBodyTextView,
		a.SendRequestButton,
		a.createHTTPResponseView(),
	)

	a.Window.SetChild(httpUIGrid)
	a.Window.SetVisible(true)

	return nil
}

func (a *View) SetOnSendRequestFunction(callback func(*http.Request)) {
	const signalClicked = "clicked"
	a.SendRequestButton.Connect(signalClicked, func() {
		urlBuffer := a.URLTextView.Buffer()
		urlStart, urlEnd := urlBuffer.Bounds()

		url := urlBuffer.Text(urlStart, urlEnd, false)

		selectedIndex := a.MethodDropDown.Selected()
		if selectedIndex == gtk.InvalidListPosition {
			a.Logger.Warn("dropdown does not have active item selected")
			return
		}

		method := a.MethodList[selectedIndex]

		bodyBuffer := a.RequestBodyTextView.Buffer()
		bodyStart, bodyEnd := bodyBuffer.Bounds()
		bodyString := bodyBuffer.Text(bodyStart, bodyEnd, true)

		var body io.ReadCloser = http.NoBody
		if len(bodyString) != 0 {
			body = io.NopCloser(strings.NewReader(bodyString))
		}

		a.Logger.Debug(
			"ready to trigger http request",
			slog.String("url", url),
			slog.String("method", method),
			slog.String("body", bodyString),
		)

		request, err := http.NewRequestWithContext(context.Background(), method, url, body)
		if err != nil {
			a.Logger.Error("cannot create request", slog.String("error", err.Error()))
			return
		}

		callback(request)
	})
}

func (a *View) createHTTPURLInput() gtk.Widgetter {
	box := gtk.NewBox(gtk.OrientationHorizontal, 0)
	box.SetSpacing(10)

	a.MethodList = []string{
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
	a.MethodDropDown = gtk.NewDropDownFromStrings(a.MethodList)
	a.MethodDropDown.SetHExpand(false)
	a.MethodDropDown.SetSelected(0)

	box.Append(a.MethodDropDown)

	a.URLTextView = gtk.NewTextView()
	a.URLTextView.SetHExpand(true)
	a.URLTextView.SetTopMargin(10)
	// TODO: make all text views rounded

	textBoxBuffer := a.URLTextView.Buffer()
	textBoxBuffer.SetText("Enter URL here")

	box.Append(a.URLTextView)

	return box
}

func (a *View) createHTTPResponseView() gtk.Widgetter {
	responseGrid := gtk.NewGrid()
	responseGrid.SetOrientation(gtk.OrientationVertical)
	responseGrid.SetRowSpacing(10)
	responseGrid.SetSizeRequest(-1, -1)
	responseGrid.SetHExpand(true)

	responseLabel := gtk.NewLabel("Response")
	responseLabel.SetHAlign(gtk.AlignStart)
	responseGrid.Attach(responseLabel, 0, 0, 1, 1)

	a.ResponseTextView = gtk.NewTextView()
	a.ResponseTextView.SetSizeRequest(-1, 100)
	a.ResponseTextView.SetEditable(false)
	a.ResponseTextView.SetHExpand(true)

	responseScrolledWindow := gtk.NewScrolledWindow()
	responseScrolledWindow.SetChild(a.ResponseTextView)
	responseScrolledWindow.SetPropagateNaturalHeight(true)

	responseGrid.Attach(responseScrolledWindow, 0, 1, 1, 1)

	return responseGrid
}

func (a *View) SetDestroyFunction(callback func()) {
	a.Window.Connect(signalDestroy, callback)
}

func (a *View) setupHTTPUserInterfaceGrid() *gtk.Grid {
	verticalGrid := gtk.NewGrid()
	verticalGrid.SetOrientation(gtk.OrientationVertical)
	verticalGrid.SetRowSpacing(10)

	verticalGrid.SetMarginStart(10)
	verticalGrid.SetMarginEnd(10)
	verticalGrid.SetMarginTop(10)
	verticalGrid.SetMarginBottom(10)

	return verticalGrid
}

func (a *View) addAllToVerticalGrid(widgets ...gtk.Widgetter) *gtk.Grid {
	grid := a.setupHTTPUserInterfaceGrid()

	for i, widget := range widgets {
		grid.Attach(widget, 0, i, 1, 1)
	}

	return grid
}
