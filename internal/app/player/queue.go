package player

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/ttpcodes/prismriver/internal/app/types"
	"sync"
)

var queueInstance *Queue
var queueOnce sync.Once

type Queue struct {
	items []types.Media
}

func GetQueue() *Queue {
	queueOnce.Do(func() {
		logrus.Info("Created queue instance.")
		queueInstance = &Queue{
			items: make([]types.Media, 0),
		}
	})
	return queueInstance
}

func (q *Queue) Add(media types.Media) {
	q.items = append(q.items, media)
	player := GetPlayer()
	if player.State == STOPPED {
		go player.Play(media)
	}
}

func (q *Queue) Advance() {
	q.items = append(q.items[:0], q.items[1:]...)
	if len(q.items) > 0 {
		player := GetPlayer()
		go player.Play(q.items[0])
	}
}

func (q Queue) GetMedia() []types.Media {
	return q.items
}
