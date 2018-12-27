package routes

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"gitlab.com/ttpcodes/prismriver/internal/app/player"
	"gitlab.com/ttpcodes/prismriver/internal/app/server/ws"
	"net/http"
	"sync"
)

var queueHub *ws.Hub
var queueOnce sync.Once

var queueUpgrader = websocket.Upgrader{
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

		go (func () {
			queue := player.GetQueue()
			for {
				select {
				case response := <- queue.Update:
					queueHub.Broadcast <- response
				}
			}
		})()
	})
	return queueHub
}

func WebsocketQueueHandler(w http.ResponseWriter, r *http.Request) {
	GetQueueHub()
	conn, err := queueUpgrader.Upgrade(w, r, nil)
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

	queue := player.GetQueue()
	response := queue.GenerateResponse()
	client.Send <- response
	logrus.Debug("Sent initial message on WS connection.")
}