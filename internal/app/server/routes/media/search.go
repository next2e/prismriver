package media

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"gitlab.com/ttpcodes/prismriver/internal/app/db"
	"net/http"
)

// SearchHandler handles requests to search for Media in the database.
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	media := db.FindMedia(query.Get("query"), 20)
	response, err := json.Marshal(media)
	if err != nil {
		logrus.Error("Error when generating SearchHandler response:")
		logrus.Error(err)
		return
	}
	w.Write(response)
}
