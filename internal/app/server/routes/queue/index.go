package queue

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"gitlab.com/ttpcodes/prismriver/internal/app/player"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	queue := player.GetQueue()
	media := queue.GetMedia()

	titles := make([]string, 0)

	for _, item := range media {
		titles = append(titles, item.Title)
	}

	response, err := json.Marshal(titles)
	if err != nil {
		logrus.Error("Error generating JSON response:")
		logrus.Error(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}