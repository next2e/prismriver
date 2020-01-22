package sources

import (
	"os"
	"path"
	"regexp"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xfrr/goffmpeg/transcoder"
	youtubedl "gitlab.com/ttpcodes/youtube-dl-go"

	"gitlab.com/ttpcodes/prismriver/internal/app/constants"
	"gitlab.com/ttpcodes/prismriver/internal/app/db"
)

// SoundCloud represents a SoundCloud Source.
type SoundCloud struct {
}

// DownloadMedia runs a download on a given item in a goroutine. This can be tracked using the returned channels.
func (s SoundCloud) DownloadMedia(media db.Media) (chan float64, chan error, error) {
	progressChan := make(chan float64)
	doneChan := make(chan error)
	go func() {
		callDone := func(err error) {
			close(progressChan)
			doneChan <- err
			close(doneChan)
		}
		downloader := youtubedl.NewDownloader("soundcloud.com/" + media.ID)

		format := viper.GetString(constants.DOWNLOADFORMAT)
		downloader.Format(format)
		downloader.Output("/tmp/" + youtubedl.ID)
		eventChan, closeChan, err := downloader.RunProgress()
		if err != nil {
			callDone(err)
			return
		}
		for progress := range eventChan {
			logrus.Debugf("download is at %f percent completion", progress)
			progressChan <- progress / 2
		}
		result := <-closeChan
		if result.Err != nil {
			callDone(result.Err)
			return
		}
		progressChan <- 50
		logrus.Debug("downloaded youtube-dl file")

		dataDir := viper.GetString(constants.DATA)
		trans := new(transcoder.Transcoder)
		filePath := path.Join(dataDir, "soundcloud", media.ID+".opus")
		if err := os.MkdirAll(path.Dir(filePath), os.ModeDir|0755); err != nil {
			callDone(err)
			return
		}
		err = trans.Initialize(result.Path, filePath)
		if err != nil {
			callDone(err)
			return
		}
		trans.MediaFile().SetAudioCodec("libopus")
		trans.MediaFile().SetSkipVideo(true)
		logrus.Debug("instantiated ffmpeg transcoder")

		done := trans.Run(true)
		progress := trans.Output()
		for msg := range progress {
			progressChan <- msg.Progress/2 + 50
			logrus.Debug(msg)
		}
		if err := <-done; err != nil {
			callDone(err)
			return
		}
		logrus.Debug("transcoded media to vorbis audio")
		if err := os.Remove(result.Path); err != nil {
			logrus.Errorf("error when removing temporary file: %v", err)
			// We don't return here because even if the temporary file isn't deleted, we successfully got the audio.
		}
		logrus.Debug("removed temporary youtube-dl file")
		logrus.Infof("downloaded new file for YouTube video ID %s", media.ID)
		callDone(nil)
	}()
	return progressChan, doneChan, nil
}

// GetInfo retrieves the info for a Media item synchronously.
func (s SoundCloud) GetInfo(id string, _ bool) (db.Media, error) {
	downloader := youtubedl.NewDownloader("soundcloud.com/" + id)
	info, err := downloader.GetInfo()
	if err != nil {
		return db.Media{}, err
	}
	return db.Media{
		ID:     id,
		Length: uint64(info.Duration * float64(time.Millisecond)),
		Title:  info.Title,
		Type:   "soundcloud",
		Video:  false,
	}, nil
}

// HasVideo indicates whether or not this Source has video capabilities.
func (s SoundCloud) HasVideo() bool {
	return false
}

// ValidateURL attempts to parse the given URL and return the Media ID, or an empty string if not parsable.
func (s SoundCloud) ValidateURL(url string) string {
	regex := regexp.MustCompile(`^https?://(soundcloud\.com|snd\.sc)/(.*)$`)
	res := regex.FindAllStringSubmatch(url, -1)
	if len(res) > 0 && len(res[0]) > 2 {
		return res[0][2]
	}
	return ""
}
