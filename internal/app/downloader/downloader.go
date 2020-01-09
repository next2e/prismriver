package downloader

import (
	"os"
	"path"

	"github.com/spf13/viper"

	"gitlab.com/ttpcodes/prismriver/internal/app/constants"
	"gitlab.com/ttpcodes/prismriver/internal/app/downloader/sources"
)

var sourceList = map[string]Source{
	"soundcloud": sources.SoundCloud{},
	"youtube":    sources.YouTube{},
}

// InitializeDownloader initializes directories and other resources needed for the download logic.
func InitializeDownloader() error {
	for name := range sourceList {
		dataDir := viper.GetString(constants.DATA)
		if err := os.MkdirAll(path.Join(dataDir, name), os.ModeDir); err != nil {
			return err
		}
	}
	return nil
}

// FindSource attempts to match the provided URL to a Source, or will return nil if not matched.
func FindSource(url string) (string, string, Source) {
	for name, source := range sourceList {
		id := source.ValidateURL(url)
		if len(id) > 0 {
			return id, name, source
		}
	}
	return "", "", nil
}

// GetSource returns a Source based on the inputted key.
func GetSource(key string) Source {
	return sourceList[key]
}
