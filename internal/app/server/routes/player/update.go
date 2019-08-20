package player

import (
	"gitlab.com/ttpcodes/prismriver/internal/app/player"
	"net/http"
)

// UpdateHandler handles requests for changes to the Player, such as calling "Be Quiet!" or modifying the volume.
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	quiet := r.Form.Get("quiet")
	if len(quiet) > 0 {
		queue := player.GetQueue()
		queue.BeQuiet()
		return
	}
	volume := r.Form.Get("volume")
	if len(volume) > 0 {
		player := player.GetPlayer()
		switch volume {
		case "up":
			player.UpVolume()
		case "down":
			player.DownVolume()
		}
	}
}
