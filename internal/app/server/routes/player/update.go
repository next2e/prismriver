package player

import (
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"

	"gitlab.com/ttpcodes/prismriver/internal/app/player"
)

// UpdateHandler handles requests for changes to the Player, such as calling "Be Quiet!" or modifying the volume.
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		logrus.Warnf("error parsing form inputs: %v", err)
		return
	}
	playerInstance := player.GetPlayer()
	queue := player.GetQueue()

	quiet := r.Form.Get("quiet")
	if len(quiet) > 0 {
		queue.BeQuiet()
	}

	seek := r.Form.Get("seek")
	if len(seek) > 0 {
		nanoseconds, err := strconv.Atoi(seek)
		if err != nil {
			logrus.Warnf("error parsing seek time: %v", err)
		} else {
			if err := playerInstance.Seek(nanoseconds); err != nil {
				logrus.Errorf("could not seek player: %v", err)
			}
		}
	}

	volume := r.Form.Get("volume")
	if len(volume) > 0 {
		playerInstance := player.GetPlayer()
		switch volume {
		case "up":
			playerInstance.UpVolume()
		case "down":
			playerInstance.DownVolume()
		}
	}
}
