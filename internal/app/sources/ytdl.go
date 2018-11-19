package sources

import (
	"gitlab.com/ttpcodes/prismriver/internal/app/constants"
	"gitlab.com/ttpcodes/prismriver/internal/app/db"
	"io/ioutil"
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

func GetVideo(query string) (string, error) {
	file, err := ioutil.TempFile("", "ytdl-")
	if err != nil {
		logrus.Error("Error when opening temporary file")
		logrus.Error(err)
		return "", err
	}
	logrus.Debug("Opened TempFile for video download")

	info, err := ytdl.GetVideoInfoFromID(query)
	if err != nil {
		logrus.Error("Error when loading video info")
		logrus.Error(err)
		return "", err
	}
	logrus.Debug("Retrieved video info")
	if err := info.Download(info.Formats.Best(ytdl.FormatAudioBitrateKey).Worst(ytdl.FormatResolutionKey)[0], file); err != nil {
		logrus.Error("Error when downloading media")
		logrus.Error(err)
		return "", err
	}
	logrus.Debug("Downloaded media file")
	if err := file.Close(); err != nil {
		logrus.Error("Error when writing temporary file")
		logrus.Error(err)
		return "", err
	}
	logrus.Debug("Wrote media to temporary file")

	trans := new(transcoder.Transcoder)
	dataDir := viper.GetString(constants.DATA)
	filePath := path.Join(dataDir, info.ID+".ogg")
	trans.Initialize(file.Name(), filePath)
	trans.MediaFile().SetAudioCodec("libvorbis")
	trans.MediaFile().SetSkipVideo(true)
	logrus.Debug("Instantiated ffmpeg transcoder")

	done := trans.Run(true)
	progress := trans.Output()
	for msg := range progress {
		logrus.Debug(msg)
	}
	if err := <-done; err != nil {
		logrus.Error("Error in transcoding process")
		logrus.Error(err)
		return "", err
	}
	logrus.Debug("Transcoded media to vorbis audio")
	if err := os.Remove(file.Name()); err != nil {
		logrus.Error("Error when removing temporary file")
		logrus.Error(err)
		return "", err
	}
	logrus.Debug("Removed temporary audio file")
	return info.ID + ".ogg", nil
}
