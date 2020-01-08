package queue

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/sirupsen/logrus"
	"gitlab.com/ttpcodes/prismriver/internal/app/db"
	"gitlab.com/ttpcodes/prismriver/internal/app/player"
	"gitlab.com/ttpcodes/prismriver/internal/app/sources"
)

// StoreHandler handles requests for adding new QueueItems.
func StoreHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		logrus.Warnf("error parsing form data from POST /queue: %v", err)
		return
	}
	id := r.Form.Get("id")
	kind := r.Form.Get("type")
	url := r.Form.Get("url")
	video, err := strconv.ParseBool(r.Form.Get("video"))
	if err != nil {
		logrus.Warnf("error parsing boolean from video input, defaulting to false")
		video = false
	}
	if len(id) > 0 && len(kind) > 0 {
		media, err := db.GetMedia(id, kind)
		if err == nil {
			queue := player.GetQueue()
			queue.Add(media)
			return
		}
	}
	if len(url) != 0 {
		regex := regexp.MustCompile(`(youtu\.be\/|youtube\.com\/(watch\?(.*&)?v=|(embed|v)\/))([^\?&"'>]+)`)
		res := regex.FindAllStringSubmatch(url, -1)
		if len(res[0]) > 0 {
			queue := player.GetQueue()
			media, err := db.GetMedia(res[0][5], "youtube")
			if err == nil {
				queue.Add(media)
				return
			}
			media, err = sources.GetInfo(res[0][5], video)
			if err != nil {
				logrus.Error("Could not get video info")
				return
			}
			if err := db.AddMedia(media); err != nil {
				logrus.Errorf("error storing new media item: %v", err)
				return
			}
			queue.Add(media)
			return
		}
	}
	logrus.Warn("User sent an empty POST request, ignoring.")
}
