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

// Download represents a download occurring for a QueueItem.
type Download struct {
	doneCh   chan struct{}
	err      string
	progress int
}

// Queue represents a queue of Media items waiting to be played.
type Queue struct {
	sync.RWMutex

	downloads map[string]*Download
	items     []*QueueItem
	Update    chan []byte
}

// QueueItem represents a Media item waiting to be played in the Queue.
type QueueItem struct {
	err   string
	Media db.Media
	ready chan bool
	queue *Queue
}

// QueueItemResponse represents a QueueItem containing the necessary fields to be exported via JSON.
type QueueItemResponse struct {
	Downloading bool     `json:"downloading"`
	Error       string   `json:"error"`
	Media       db.Media `json:"media"`
	Progress    int      `json:"progress"`
}

// GetQueue returns the single Queue instance of the application.
func GetQueue() *Queue {
	queueOnce.Do(func() {
		logrus.Info("Created queue instance.")
		queueInstance = &Queue{
			downloads: make(map[string]*Download),
			items:     make([]*QueueItem, 0),
			Update:    make(chan []byte),
		}
	})
	return queueInstance
}

// Add adds a new Media item to the Queue as a QueueItem. If the item is detected to not be ready, it will instantiate
// a download of the Media.
func (q *Queue) Add(media db.Media) {
	item := &QueueItem{
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
		// Read of downloads and assignment to map must be done together or another lock might get hold of it.
		q.Lock()
		if download, ok := q.downloads[item.Media.ID]; ok {
			q.Unlock()
			go func() {
				<-download.doneCh
				if download.err != "" {
					item.err = download.err
					q.sendQueueUpdate()
					return
				}
				item.ready <- true
				close(item.ready)
			}()
		} else {
			download := &Download{
				doneCh: make(chan struct{}),
			}
			q.downloads[item.Media.ID] = download
			q.Unlock()

			q.sendQueueUpdate()
			progressChan, doneChan, err := sources.GetVideo(media.ID)
			if err != nil {
				logrus.Error("Error when getting sources:")
				logrus.Error(err)
				return
			}

			go func() {
				for progress := range progressChan {
					q.Lock()
					download.progress = int(progress)
					q.Unlock()
					q.sendQueueUpdate()
				}
				if err := <-doneChan; err != nil {
					q.Lock()
					download.err = err.Error()
					item.err = err.Error()
					delete(q.downloads, item.Media.ID)
					q.Unlock()
					close(download.doneCh)
					q.sendQueueUpdate()
					return
				}
				q.Lock()
				delete(q.downloads, item.Media.ID)
				q.Unlock()
				close(download.doneCh)
				q.sendQueueUpdate()
				item.ready <- true
				close(item.ready)
			}()
		}
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

// Advance moves the Queue up by one and plays the next item if it exists.
func (q *Queue) Advance() {
	q.items = q.items[1:]
	if len(q.items) > 0 {
		player := GetPlayer()
		go player.Play(q.items[0])
	}
	go q.sendQueueUpdate()
}

// BeQuiet replaces the currently playing item with the BeQuiet Media and plays it.
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

// MoveDown moves a QueueItem down in the Queue.
func (q *Queue) MoveDown(index int) {
	if index == len(q.items)-1 {
		return
	}
	temp := q.items[index+1]
	q.items[index+1] = q.items[index]
	q.items[index] = temp
	q.sendQueueUpdate()
}

// MoveUp moves a QueueItem up in the Queue.
func (q *Queue) MoveUp(index int) {
	if index == 1 {
		return
	}
	temp := q.items[index-1]
	q.items[index-1] = q.items[index]
	q.items[index] = temp
	q.sendQueueUpdate()
}

// GenerateResponse generates a JSON response of all the QueueItems in the Queue.
func (q *Queue) GenerateResponse() []byte {
	// Cannot return a nil slice or the frontend will have issues.
	items := make([]QueueItemResponse, 0)
	q.RLock()
	for _, item := range q.items {
		response := item.GenerateResponse()
		items = append(items, response)
	}
	q.RUnlock()
	response, err := json.Marshal(items)
	if err != nil {
		logrus.Error("Error generating JSON response:")
		logrus.Error(err)
	}
	return response
}

// GetItems returns the QueueItems in the Queue.
func (q *Queue) GetItems() []*QueueItem {
	return q.items
}

// Remove removes a QueueItem from the Queue.
func (q *Queue) Remove(index int) {
	q.items = append(q.items[:index], q.items[index+1:]...)
	go q.sendQueueUpdate()
}

func (q *Queue) sendQueueUpdate() {
	response := q.GenerateResponse()
	q.Update <- response
}

// GenerateResponse returns the QueueItemResponse form of the QueueItem.
func (q QueueItem) GenerateResponse() QueueItemResponse {
	downloading, progress := q.Progress()
	return QueueItemResponse{
		Downloading: downloading,
		Error:       q.err,
		Media:       q.Media,
		Progress:    progress,
	}
}

// Progress returns the download progress of the QueueItem.
func (q QueueItem) Progress() (bool, int) {
	download, ok := q.queue.downloads[q.Media.ID]
	if !ok {
		return false, 100
	}
	return true, download.progress
}
