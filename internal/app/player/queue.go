package player

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"gitlab.com/ttpcodes/prismriver/internal/app/db"
	"gitlab.com/ttpcodes/prismriver/internal/app/server/ws/routes"
	"sync"
)

var queueInstance *Queue
var queueOnce sync.Once

type Queue struct {
	items []db.Media
}

func GetQueue() *Queue {
	queueOnce.Do(func() {
		logrus.Info("Created queue instance.")
		queueInstance = &Queue{
			items: make([]db.Media, 0),
		}
	})
	return queueInstance
}

func (q *Queue) Add(media db.Media) {
	q.items = append(q.items, media)
	player := GetPlayer()
	if player.State == STOPPED {
		go player.Play(media)
	}
	go q.sendQueueUpdate()
	logrus.Info("Added " + media.Title + " to queue.")
}

func (q *Queue) Advance() {
	q.items = append(q.items[:0], q.items[1:]...)
	if len(q.items) > 0 {
		player := GetPlayer()
		go player.Play(q.items[0])
	}
	go q.sendQueueUpdate()
}

func (q Queue) GetMedia() []db.Media {
	return q.items
}

func (q Queue) sendQueueUpdate() {
	logrus.Debug("Called sendQueueUpdate.")
	titles := make([]string, 0)
	for _, item := range q.items {
		titles = append(titles, item.Title)
	}
	message, err := json.Marshal(titles)
	if err != nil {
		logrus.Error("Error generating JSON response:")
		logrus.Error(err)
	}
	hub := routes.GetQueueHub()
	hub.Broadcast <- message
}