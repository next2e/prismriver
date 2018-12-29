package queue

import (
	"gitlab.com/ttpcodes/prismriver/internal/app/player"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	queue := player.GetQueue()
	response := queue.GenerateResponse()
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}