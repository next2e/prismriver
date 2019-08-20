package queue

import (
	"gitlab.com/ttpcodes/prismriver/internal/app/player"
	"net/http"
)

// IndexHandler handles requests for listing all QueueItems.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	queue := player.GetQueue()
	response := queue.GenerateResponse()
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
