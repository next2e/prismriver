package player

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"gitlab.com/ttpcodes/prismriver/internal/app/db"
	"sync"
)

var queueInstance *Queue
var queueOnce sync.Once

type Queue struct {
	items []db.Media
	Update chan []byte
}

func GetQueue() *Queue {
	queueOnce.Do(func() {
		logrus.Info("Created queue instance.")
		queueInstance = &Queue{
			items: make([]db.Media, 0),
			Update: make(chan []byte),
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

func (q *Queue) MoveDown(index int) {
	if index == len(q.items) - 1 {
		return
	}
	temp := q.items[index + 1]
	q.items[index + 1] = q.items[index]
	q.items[index] = temp
	q.sendQueueUpdate()
}

func (q *Queue) MoveUp(index int) {
	if index == 1 {
		return
	}
	temp := q.items[index - 1]
	q.items[index - 1] = q.items[index]
	q.items[index] = temp
	q.sendQueueUpdate()
}

func (q Queue) GenerateResponse() []byte {
	titles := make([]string, 0)
	for _, item := range q.items {
		titles = append(titles, item.Title)
	}
	response, err := json.Marshal(titles)
	if err != nil {
		logrus.Error("Error generating JSON response:")
		logrus.Error(err)
	}
	return response
}

func (q Queue) GetMedia() []db.Media {
	return q.items
}

func (q *Queue) Remove(index int) {
	q.items = append(q.items[:index], q.items[index + 1:]...)
	go q.sendQueueUpdate()
}

func (q Queue) sendQueueUpdate() {
	response := q.GenerateResponse()
	q.Update <- response
}