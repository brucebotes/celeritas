package websocket

import (
	"fmt"

	"github.com/pusher/pusher-http-go"
)

type WebSocket struct {
	Secret       string
	Key          string
	Host         string
	Port         string
	AuthEndPoint string
	Secure       bool
}

var wsClient pusher.Client

func (w *WebSocket) Init(appID string) *pusher.Client {
	// create pusher client
	wsClient = pusher.Client{
		AppID:  appID,
		Secret: w.Secret,
		Key:    w.Key,
		Secure: w.Secure,
		Host:   fmt.Sprintf("%s:%s", w.Host, w.Port),
	}

	return &wsClient
}
