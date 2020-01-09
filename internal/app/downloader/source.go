package downloader

import "gitlab.com/ttpcodes/prismriver/internal/app/db"

// Source is an interface representing a media source that can be downloaded from.
type Source interface {
	DownloadMedia(media db.Media) (chan float64, chan error, error)
	GetInfo(id string, video bool) (db.Media, error)
	HasVideo() bool
	ValidateURL(url string) string
}
