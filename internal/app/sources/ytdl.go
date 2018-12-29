package sources

import (
	"gitlab.com/ttpcodes/prismriver/internal/app/constants"
	"gitlab.com/ttpcodes/prismriver/internal/app/db"
	"gitlab.com/ttpcodes/youtube-dl-go"
	"os"
	"path"

	"github.com/rylio/ytdl"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xfrr/goffmpeg/transcoder"
)

func GetInfo(query string) (db.Media, error) {
	info, err := ytdl.GetVideoInfo(query)
	if err != nil {
		logrus.Error("Error retrieving video info:")
		logrus.Error(err)
		return db.Media{}, err
	}
	return db.Media{
		ID: info.ID,
		Length: uint64(info.Duration),
		Title: info.Title,
		Type: "youtube",
	}, nil
}

func GetVideo(query string) (chan float64, chan struct{}, error) {
	info, err := ytdl.GetVideoInfoFromID(query)
	if err != nil {
		logrus.Error("Error when loading video info")
		logrus.Error(err)
		return nil, nil, err
	}
	logrus.Debug("Retrieved video info")
	progressChan := make(chan float64)
	doneChan := make(chan struct{})
	go func() {
		downloader := youtubedl.NewDownloader(query)
		downloader.Output("/tmp/" + youtubedl.ID)
		eventChan, closeChan, err := downloader.RunProgress()
		if err != nil {
			logrus.Error("Error when downloading video file:")
			logrus.Error(err)
		}
		for progress := range eventChan {
			logrus.Debugf("Download is at %f percent completion", progress)
			progressChan <- progress / 2
		}
		tmpPath := <- closeChan
		logrus.Debug("Downloaded media file")

		trans := new(transcoder.Transcoder)
		dataDir := viper.GetString(constants.DATA)
		filePath := path.Join(dataDir, info.ID+".opus")
		trans.Initialize(tmpPath, filePath)
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
			logrus.Error("Error in transcoding process")
			logrus.Error(err)
		}
		logrus.Debug("Transcoded media to vorbis audio")
		if err := os.Remove(tmpPath); err != nil {
			logrus.Error("Error when removing temporary file")
			logrus.Error(err)
		}
		logrus.Debug("Removed temporary audio file")
		logrus.Infof("Downloaded new audio file for YouTube video ID %s", info.ID)
		close(progressChan)
		doneChan <- struct{}{}
		close(doneChan)
	}()
	return progressChan, doneChan, nil
}
