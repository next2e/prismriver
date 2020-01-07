package sources

import (
	"gitlab.com/ttpcodes/prismriver/internal/app/constants"
	"gitlab.com/ttpcodes/prismriver/internal/app/db"
	"gitlab.com/ttpcodes/youtube-dl-go"
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xfrr/goffmpeg/transcoder"
)

// GetInfo retrieves the info for the YouTube video and returns it as a Media item, or an error if encountered.
func GetInfo(id string) (db.Media, error) {
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
	}, nil
}

// GetVideo attempts to download a YouTube video specified by id. It will return channels that can be used to track
// the progress of the download and any errors encountered when the process finishes.
func GetVideo(id string) (chan float64, chan error, error) {
	progressChan := make(chan float64)
	doneChan := make(chan error)
	go func() {
		callDone := func(err error) {
			close(progressChan)
			doneChan <- err
			close(doneChan)
		}
		downloader := youtubedl.NewDownloader(id)
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

		trans := new(transcoder.Transcoder)
		dataDir := viper.GetString(constants.DATA)
		filePath := path.Join(dataDir, id+".opus")
		err = trans.Initialize(result.Path, filePath)
		if err != nil {
			logrus.Error("Error starting transcoding process:\n", err)
			callDone(err)
			return
		}
		trans.MediaFile().SetAudioCodec("libopus")
		trans.MediaFile().SetSkipVideo(true)
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
		if err := os.Remove(result.Path); err != nil {
			logrus.Error("Error when removing temporary file")
			logrus.Error(err)
			// We don't return here because even if the temporary file isn't deleted, we successfully got the audio.
		}
		logrus.Debug("Removed temporary audio file")
		logrus.Infof("Downloaded new audio file for YouTube video ID %s", id)
		callDone(nil)
	}()
	return progressChan, doneChan, nil
}
