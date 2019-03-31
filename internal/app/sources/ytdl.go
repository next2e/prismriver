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

func GetInfo(query string) (db.Media, error) {
	downloader := youtubedl.NewDownloader(query)
	info, err := downloader.GetInfo()
	if err != nil {
		logrus.Error("Error retrieving video info:")
		logrus.Error(err)
		return db.Media{}, err
	}
	return db.Media{
		ID: info.ID,
		Length: uint64(info.Duration * 1000000),
		Title: info.Title,
		Type: "youtube",
	}, nil
}

func GetVideo(query string) (chan float64, chan error, error) {
	downloader := youtubedl.NewDownloader(query)
	info, err := downloader.GetInfo()
	if err != nil {
		logrus.Error("Error when loading video info")
		logrus.Error(err)
		return nil, nil, err
	}
	logrus.Debug("Retrieved video info")
	progressChan := make(chan float64)
	doneChan := make(chan error)
	go func() {
		downloader := youtubedl.NewDownloader(query)
		downloader.Output("/tmp/" + youtubedl.ID)
		eventChan, closeChan, err := downloader.RunProgress()
		if err != nil {
			logrus.Error("Error when downloading video file:\n", err)
			doneChan <- err
			close(doneChan)
			return
		}
		for progress := range eventChan {
			logrus.Debugf("Download is at %f percent completion", progress)
			progressChan <- progress / 2
		}
		result := <- closeChan
		if result.Err != nil {
			logrus.Error("Error downloading media file:\n", err)
			doneChan <- err
			close(doneChan)
			return
		}
		logrus.Debug("Downloaded media file")

		trans := new(transcoder.Transcoder)
		dataDir := viper.GetString(constants.DATA)
		filePath := path.Join(dataDir, info.ID+".opus")
		err = trans.Initialize(result.Path, filePath)
		if err != nil {
			logrus.Error("Error starting transcoding process:\n", err)
			doneChan <- err
			close(doneChan)
			return
		}
		trans.MediaFile().SetAudioCodec("libopus")
		trans.MediaFile().SetSkipVideo(true)
		logrus.Debug("Instantiated ffmpeg transcoder")

		done := trans.Run(true)
		progress := trans.Output()
		for msg := range progress {
			progressChan <- msg.Progress / 2 + 50
			logrus.Debug(msg)
		}
		if err := <-done; err != nil {
			logrus.Error("Error in transcoding process:\n", err)
			doneChan <- err
			close(doneChan)
			return
		}
		logrus.Debug("Transcoded media to vorbis audio")
		if err := os.Remove(result.Path); err != nil {
			logrus.Error("Error when removing temporary file")
			logrus.Error(err)
			// We don't return here because even if the temporary file isn't deleted, we successfully got the audio.
		}
		logrus.Debug("Removed temporary audio file")
		logrus.Infof("Downloaded new audio file for YouTube video ID %s", info.ID)
		close(progressChan)
		doneChan <- nil
		close(doneChan)
	}()
	return progressChan, doneChan, nil
}
