package controller

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/quic-go/quic-go/http3"
	"github.com/vin-rmdn/http3-gui-client-go/internal/view"
)

type controller struct {
	httpClient *http.Client
	view       *view.View
}

func New(httpClient *http.Client, view *view.View) controller {
	c := controller{
		httpClient: httpClient,
		view:       view,
	}

	view.SetOnSendRequestFunction(c.handleSendRequestButton)
	view.SetDestroyFunction(c.Destroy)

	return c
}

func (c controller) handleSendRequestButton(request *http.Request) {
	response, err := c.httpClient.Do(request)
	if err != nil {
		slog.Error("cannot execute http call", slog.String("error", err.Error()))
		return
	}

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

	c.view.SetResponseBody(string(responseBody))
}

func (c controller) Destroy() {
	c.httpClient.Transport.(*http3.Transport).Close()
}
