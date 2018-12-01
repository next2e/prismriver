package player

import (
	"github.com/adrg/libvlc-go"
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
	LOADING = iota
	PAUSED  = iota
	PLAYING = iota
	STOPPED = iota
)

type Player struct {
	doneChan chan bool
	player   *vlc.Player
	State    int
	Volume   int
}

func GetPlayer() *Player {
	playerOnce.Do(func() {
		playerInstance = &Player{
			doneChan: make(chan bool),
			State:    STOPPED,
			Volume:   100,
		}
	})
	return playerInstance
}

func (p *Player) Play(media db.Media) error {
	defer func() {
		p.State = STOPPED
		queue := GetQueue()
		queue.Advance()
	}()
	p.State = LOADING
	dataDir := viper.GetString(constants.DATA)
	filePath := path.Join(dataDir, media.ID+".opus")

	defer vlc.Release()

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		sources.GetVideo(media.ID)
	}

	if err := vlc.Init("--no-video", "--quiet"); err != nil {
		logrus.Error("Error initializing vlc:")
		logrus.Error(err)
		return err
	}
	defer vlc.Release()

	var err error
	p.player, err = vlc.NewPlayer()
	if err != nil {
		logrus.Error("Error creating player:")
		logrus.Error(err)
		return err
	}
	defer func() {
		p.player.Stop()
		p.player.Release()
	}()

	vlcMedia, err := p.player.LoadMediaFromPath(filePath)
	if err != nil {
		logrus.Error("Error loading media file:")
		logrus.Error(err)
		return err
	}
	defer vlcMedia.Release()

	if err := p.player.Play(); err != nil {
		logrus.Error("Error playing media file:")
		logrus.Error(err)
		return err
	}
	p.player.SetVolume(p.Volume)

	time.Sleep(1 * time.Second)
	length, err := p.player.MediaLength()
	if err != nil || length == 0 {
		length = 1000 * 60
	}
	p.State = PLAYING

	select {
	case <-p.doneChan:
		break
	case <-time.After(time.Duration(length) * time.Millisecond):
		break
	}

	return nil
}

func (p *Player) Skip() {
	logrus.Debug("")
	p.doneChan <- true
}

func (p *Player) UpVolume() {
	if p.Volume == 100 {
		return
	}
	p.Volume += 5
	p.player.SetVolume(p.Volume)
}

func (p *Player) DownVolume() {
	if p.Volume == 0 {
		return
	}
	p.Volume -= 5
	p.player.SetVolume(p.Volume)
}
