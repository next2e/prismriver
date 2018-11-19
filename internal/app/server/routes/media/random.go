package media

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"gitlab.com/ttpcodes/prismriver/internal/app/db"
	"net/http"
	"strconv"
)

func RandomHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	str := query.Get("limit")
	limit, err := strconv.ParseInt(str, 10, 8)
	if err != nil {
		logrus.Warn("Could not properly parse limit on RandomHandler.")
	}
	media := db.GetRandomMedia(int(limit))
	response, err := json.Marshal(media)
	if err != nil {
		logrus.Error("Error when generating RandomHandler response:")
		logrus.Error(err)
		return
	}
	w.Write(response)
}
