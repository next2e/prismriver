package sources

import (
	"io"
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

// YouTube represents a YouTube Source.
type YouTube struct {
}

// DownloadMedia runs a download on a given item in a goroutine. This can be tracked using the returned channels.
func (y YouTube) DownloadMedia(media db.Media) (chan float64, chan error, error) {
	progressChan := make(chan float64)
	doneChan := make(chan error)
	go func() {
		callDone := func(err error) {
			close(progressChan)
			doneChan <- err
			close(doneChan)
		}
		downloader := youtubedl.NewDownloader(media.ID)

		format := viper.GetString(constants.DOWNLOADFORMAT)
		downloader.Format(format)
		downloader.Output("/tmp/" + youtubedl.ID)
		eventChan, closeChan, err := downloader.RunProgress()
		if err != nil {
			logrus.Error("Error starting media download:\n", err)
			callDone(err)
			return
		}
		for progress := range eventChan {
			logrus.Debugf("Download is at %f percent completion", progress)
			progressChan <- progress / 2
		}
		result := <-closeChan
		if result.Err != nil {
			logrus.Error("Error downloading media file:\n", result.Err)
			callDone(result.Err)
			return
		}
		logrus.Debug("Downloaded media file")

		dataDir := viper.GetString(constants.DATA)
		if !media.Video || viper.GetBool(constants.VIDEOTRANSCODING) {
			trans := new(transcoder.Transcoder)
			ext := ".opus"
			if media.Video {
				ext = ".mp4"
			}
			filePath := path.Join(dataDir, "youtube", media.ID+ext)
			err = trans.Initialize(result.Path, filePath)
			if err != nil {
				logrus.Error("Error starting transcoding process:\n", err)
				callDone(err)
				return
			}
			trans.MediaFile().SetAudioCodec("libopus")
			if media.Video {
				trans.MediaFile().SetVideoCodec("libx264")
				// Needed to enable experimental Opus in the mp4 container format.
				trans.MediaFile().SetStrict(-2)
			} else {
				trans.MediaFile().SetSkipVideo(true)
			}
			logrus.Debug("Instantiated ffmpeg transcoder")

			done := trans.Run(true)
			progress := trans.Output()
			for msg := range progress {
				progressChan <- msg.Progress/2 + 50
				logrus.Debug(msg)
			}
			if err := <-done; err != nil {
				logrus.Error("Error in transcoding process:\n", err)
				callDone(err)
				return
			}
			logrus.Debug("Transcoded media to vorbis audio")
		} else {
			logrus.Debugf("video transcoding disabled, moving file to final destination")
			input, err := os.Open(result.Path)
			if err != nil {
				logrus.Errorf("error reading original video file: %v", err)
				callDone(err)
				return
			}
			defer func() {
				if err := input.Close(); err != nil {
					logrus.Errorf("error closing input file: %v", err)
				}
			}()
			output, err := os.Create(path.Join(dataDir, "youtube", media.ID+".video"))
			if err != nil {
				logrus.Errorf("error opening destination file: %v", err)
				callDone(err)
				return
			}
			defer func() {
				if err := output.Close(); err != nil {
					logrus.Errorf("error closing output file: %v", err)
				}
			}()
			if _, err := io.Copy(output, input); err != nil {
				logrus.Errorf("error copying video file: %v", err)
				callDone(err)
				return
			}
		}
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
func (y YouTube) GetInfo(id string, video bool) (db.Media, error) {
	downloader := youtubedl.NewDownloader(id)
	info, err := downloader.GetInfo()
	if err != nil {
		return db.Media{}, err
	}
	return db.Media{
		ID:     id,
		Length: uint64(info.Duration * float64(time.Millisecond)),
		Title:  info.Title,
		Type:   "youtube",
		Video:  video,
	}, nil
}

// HasVideo indicates whether or not this Source has video capabilities.
func (y YouTube) HasVideo() bool {
	return true
}

// ValidateURL attempts to parse the given URL and return the Media ID, or an empty string if not parsable.
func (y YouTube) ValidateURL(url string) string {
	regex := regexp.MustCompile(`(youtu\.be/|youtube\.com/(watch\?(.*&)?v=|(embed|v)/))([^?&"'>]+)`)
	res := regex.FindAllStringSubmatch(url, -1)
	if len(res) > 0 && len(res[0]) > 5 {
		return res[0][5]
	}
	return ""
}
