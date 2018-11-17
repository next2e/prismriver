package queue

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/ttpcodes/prismriver/internal/app/player"
	"gitlab.com/ttpcodes/prismriver/internal/app/sources"
	"net/http"
)

func StoreHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	url := r.Form.Get("url")
	if len(url) == 0 {
		logrus.Warn("User sent an empty POST request, ignoring.")
	} else {
		queue := player.GetQueue()
		media, err := sources.GetInfo(url)
		if err != nil {
			logrus.Error("Could not get video info")
			return
		}
		queue.Add(media)
		logrus.Info("Added " + media.Title + " to queue.")
	}

}
