package player

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ttpcodes/prismriver/internal/app/constants"
	"gitlab.com/ttpcodes/prismriver/internal/app/db"
	"gitlab.com/ttpcodes/prismriver/internal/app/sources"
	"os"
	"path"
	"sync"
)

var queueInstance *Queue
var queueOnce sync.Once

type Queue struct {
	items []*QueueItem
	Update chan []byte
}

type QueueItem struct {
	Downloading bool
	DownloadProgress float64
	Media db.Media
	ready chan bool
	queue *Queue
}

func GetQueue() *Queue {
	queueOnce.Do(func() {
		logrus.Info("Created queue instance.")
		queueInstance = &Queue{
			items: make([]*QueueItem, 0),
			Update: make(chan []byte),
		}
	})
	return queueInstance
}

func (q *Queue) Add(media db.Media) {
	item := &QueueItem{
		Downloading: false,
		Media: media,
		ready: make(chan bool),
		queue: q,
	}
	q.items = append(q.items, item)
	length := len(q.items)
	q.sendQueueUpdate()

	dataDir := viper.GetString(constants.DATA)
	filePath := path.Join(dataDir, media.ID+".opus")
	_, err := os.Stat(filePath)
	if item.Media.Type != "internal" && os.IsNotExist(err) {
		item.UpdateDownload(true, 0)
		progressChan, doneChan, err := sources.GetVideo(media.ID)
		if err != nil {
			logrus.Error("Error when getting sources:")
			logrus.Error(err)
			return
		}

		go func() {
			for progress := range progressChan {
				item.UpdateDownload(true, progress)
			}
			<- doneChan
			item.UpdateDownload(false, 100)
			item.ready <- true
			close(item.ready)
		}()
	} else {
		logrus.Debug("Queue item ready. Sending on channel.")
		go func() {
			item.ready <- true
			close(item.ready)
		}()
	}
	player := GetPlayer()
	if player.State == STOPPED && length == 1 {
		go player.Play(item)
	}
	q.sendQueueUpdate()
	logrus.Info("Added " + media.Title + " to queue.")
}

func (q *Queue) Advance() {
	q.items = q.items[1:]
	if len(q.items) > 0 {
		player := GetPlayer()
		go player.Play(q.items[0])
	}
	go q.sendQueueUpdate()
}

func (q *Queue) BeQuiet() {
	player := GetPlayer()
	if len(q.items) == 0 {
		q.Add(*db.BeQuiet)
		return
	} else if player.State == LOADING {
		return
	}
	quietQueue := make([]*QueueItem, 0)
	quietItem := &QueueItem{
		Downloading: false,
		Media: *db.BeQuiet,
		ready: make(chan bool, 1),
		queue: q,
	}
	quietItem.ready <- true
	close(quietItem.ready)
	quietQueue = append(quietQueue, q.items[0], quietItem)
	quietQueue = append(quietQueue, q.items[1:]...)
	q.items = quietQueue
	player.Skip()
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
	response, err := json.Marshal(q.items)
	if err != nil {
		logrus.Error("Error generating JSON response:")
		logrus.Error(err)
	}
	return response
}

func (q Queue) GetItems() []*QueueItem {
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

func (q *QueueItem) UpdateDownload(downloading bool, progress float64) {
	q.Downloading = downloading
	q.DownloadProgress = progress
	go q.queue.sendQueueUpdate()
}