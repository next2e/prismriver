package player

import (
	"encoding/json"
	"github.com/adrg/libvlc-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ttpcodes/prismriver/internal/app/constants"
	"path"
	"sync"
	"time"
)

var playerInstance *Player
var playerOnce sync.Once
var playerTicker *time.Ticker

const (
	STOPPED = iota
	PLAYING = iota
	PAUSED = iota
	LOADING = iota
)

type Player struct {
	doneChan chan bool
	player   *vlc.Player
	State    int
	Update   chan []byte
	Volume   int
}

type PlayerState struct {
	CurrentTime int
	TotalTime int
	State int
	Volume int
}

func GetPlayer() *Player {
	playerOnce.Do(func() {
		playerInstance = &Player{
			doneChan: make(chan bool),
			State:    STOPPED,
			Update:   make(chan []byte),
			Volume:   100,
		}
		playerTicker = time.NewTicker(30 * time.Second)
		go func () {
			for {
				select {
				case <- playerTicker.C:
					response := playerInstance.GenerateResponse()
					playerInstance.Update <- response
				}
			}
		}()
	})
	return playerInstance
}

func (p Player) GenerateResponse() []byte {
	if p.State == PLAYING {
		currentTime, err := p.player.MediaTime()
		if err != nil {
			logrus.Error("Error getting player's media time:")
			logrus.Error(err)
		}
		totalTime, err := p.player.MediaLength()
		if err != nil {
			logrus.Error("Error getting player's media length:")
			logrus.Error(err)
		}
		response, err := json.Marshal(PlayerState{
			CurrentTime: currentTime,
			State: p.State,
			TotalTime: totalTime,
			Volume: p.Volume,
		})
		if err != nil {
			logrus.Error("Error generating JSON response:")
			logrus.Error(err)
		}
		return response
	} else {
		response, err := json.Marshal(PlayerState{
			CurrentTime: 0,
			State: p.State,
			TotalTime: 0,
			Volume: p.Volume,
		})
		if err != nil {
			logrus.Error("Error generating JSON response:")
			logrus.Error(err)
		}
		return response
	}

}

func (p *Player) Play(item *QueueItem) error {
	defer func() {
		p.State = STOPPED
		queue := GetQueue()
		queue.Advance()
	}()
	p.State = LOADING
	dataDir := viper.GetString(constants.DATA)
	var filePath string
	if item.Media.Type == "internal" {
		filePath = path.Join(dataDir, "internal", item.Media.ID+".opus")
	} else {
		filePath = path.Join(dataDir, item.Media.ID+".opus")
	}
	ready := <- item.ready

	if !ready {
		logrus.Warn("Item labeled as not ready, not playing.")
		return nil
	}

	defer vlc.Release()

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
	p.sendPlayerUpdate()

	select {
	case <-p.doneChan:
		break
	case <-time.After(time.Duration(length) * time.Millisecond):
		break
	}
	p.State = STOPPED
	p.sendPlayerUpdate()

	return nil
}

func (p *Player) Skip() {
	logrus.Debug("")
	if p.State == PLAYING {
		p.doneChan <- true
	} else {
		queue := GetQueue()
		queue.items[0].ready <- false
	}
}

func (p *Player) UpVolume() {
	if p.Volume == 100 {
		return
	}
	p.Volume += 5
	if p.State == PLAYING {
		p.player.SetVolume(p.Volume)
	}
	p.sendPlayerUpdate()
}

func (p *Player) DownVolume() {
	if p.Volume == 0 {
		return
	}
	p.Volume -= 5
	if p.State == PLAYING {
		p.player.SetVolume(p.Volume)
	}
	p.sendPlayerUpdate()
}

func (p Player) sendPlayerUpdate() {
	response := p.GenerateResponse()
	p.Update <- response
	logrus.Debug("Sent player update event.")
}