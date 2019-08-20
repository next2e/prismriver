package queue

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gitlab.com/ttpcodes/prismriver/internal/app/player"
	"net/http"
	"strconv"
)

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	index, err := strconv.ParseUint(vars["id"], 10, 8)
	if err != nil {
		logrus.Warn("Error parsing int in DeleteHandler, user likely provided incorrect input.")
		return
	}
	if int(index) == 0 {
		playInstance := player.GetPlayer()
		playInstance.Skip()
	} else {
		queue := player.GetQueue()
		items := queue.GetItems()

		if int(index) < len(items) {
			queue.Remove(int(index))
		}
	}
}
