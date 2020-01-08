package sources

import (
	"io"
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xfrr/goffmpeg/transcoder"

	"gitlab.com/ttpcodes/prismriver/internal/app/constants"
	"gitlab.com/ttpcodes/prismriver/internal/app/db"
	"gitlab.com/ttpcodes/youtube-dl-go"
)

// GetInfo retrieves the info for the YouTube video and returns it as a Media item, or an error if encountered.
func GetInfo(id string, video bool) (db.Media, error) {
	downloader := youtubedl.NewDownloader(id)
	info, err := downloader.GetInfo()
	if err != nil {
		logrus.Error("Error retrieving video info:")
		logrus.Error(err)
		return db.Media{}, err
	}
	return db.Media{
		ID:     info.ID,
		Length: uint64(info.Duration * 1000000),
		Title:  info.Title,
		Type:   "youtube",
		Video:  video,
	}, nil
}

// GetVideo attempts to download a YouTube video specified by id. It will return channels that can be used to track
// the progress of the download and any errors encountered when the process finishes.
func GetVideo(media db.Media) (chan float64, chan error, error) {
	progressChan := make(chan float64)
	doneChan := make(chan error)
	go func() {
		callDone := func(err error) {
			close(progressChan)
			doneChan <- err
			close(doneChan)
		}
		downloader := youtubedl.NewDownloader(media.ID)
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
			filePath := path.Join(dataDir, media.ID+ext)
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
			output, err := os.Create(path.Join(dataDir, media.ID+".video"))
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
