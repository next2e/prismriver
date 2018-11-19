package player

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ttpcodes/prismriver/internal/app/constants"
	"gitlab.com/ttpcodes/prismriver/internal/app/db"
	"gitlab.com/ttpcodes/prismriver/internal/app/sources"
	"os"
	"path"
	"sync"
	"time"
)

var playerInstance *Player
var playerOnce sync.Once

const (
	PAUSED  = iota
	PLAYING = iota
	STOPPED = iota
)

type Player struct {
	State int
}

func GetPlayer() *Player {
	playerOnce.Do(func() {
		playerInstance = &Player{
			State: STOPPED,
		}
	})
	return playerInstance
}

func (p *Player) Play(media db.Media) error {
	dataDir := viper.GetString(constants.DATA)
	filePath := path.Join(dataDir, media.ID+".ogg")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		sources.GetVideo(media.ID)
	}

	file, err := os.Open(filePath)
	if err != nil {
		logrus.Error("Error opening media file:")
		logrus.Error(err)
		return err
	}
	stream, format, err := vorbis.Decode(file)
	if err != nil {
		logrus.Error("Error decoding media file:")
		logrus.Error(err)
		return err
	}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan struct{})
	speaker.Play(beep.Seq(stream, beep.Callback(func() {
		close(done)
	})))
	playerInstance.State = PLAYING
	<-done
	playerInstance.State = STOPPED
	queue := GetQueue()
	queue.Advance()
	return nil
}
