package routes

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"gitlab.com/ttpcodes/prismriver/internal/app/server/ws"
	"net/http"
	"sync"
)

var queueHub *ws.Hub
var queueOnce sync.Once

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

func GetQueueHub() *ws.Hub {
	queueOnce.Do(func () {
		queueHub = ws.CreateHub()
		go queueHub.Execute()
	})
	return queueHub
}

func WebsocketQueueHandler(w http.ResponseWriter, r *http.Request) {
	GetQueueHub()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Error("Error when upgrading client to WS connection:")
		logrus.Error(err)
	}
	client := &ws.Client{
		Conn: conn,
		Hub: queueHub,
		Send: make(chan []byte, 256),
	}
	client.Hub.Register <- client

	go client.RunRead()
	go client.RunWrite()
}