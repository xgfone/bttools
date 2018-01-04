package downloader

import (
	"errors"
	"sync"
	"time"

	"github.com/xgfone/go-tools/net2"
)

var (
	// ErrNotFound is returned when cannot found the .torrent file.
	ErrNotFound = errors.New("Not Found")

	// ErrServerRefuse is returned when the server refuse to connect to it.
	ErrServerRefuse = errors.New("Server Refuse Error")

	// ErrInvalidHash is returned when the infohash is not right.
	ErrInvalidHash = errors.New("invalid infohash")

	// ErrExist is returned when an item has existed.
	ErrExist = errors.New("Has existed")

	// ErrNotExist is returned when an item has not existed.
	ErrNotExist = errors.New("Has not existed")

	// ErrEmptyURL is returned when url is empty.
	ErrEmptyURL = errors.New("Url is empty")
)

var (
	downloadersLock = new(sync.Mutex)
	downloaders     = make(map[string]Downloader)

	downlandedAddress = ""
)

// SetHostPost resets the host and the port of the downloaded url.
func SetHostPost(host string, port int) {
	downlandedAddress = net2.JoinHostPort(host, port)
}

func getDownloadedAddress(_defaultAddr string) string {
	if downlandedAddress != "" {
		return downlandedAddress
	}
	return _defaultAddr
}

// RegisterDownloader registers a downloader.
func RegisterDownloader(name string, d Downloader) (err error) {
	downloadersLock.Lock()
	if _, ok := downloaders[name]; ok {
		err = ErrExist
	} else {
		downloaders[name] = d
	}
	downloadersLock.Unlock()
	return
}

// GetDownloader returns a downloader named name.
// Return nil if there is no downloader named name.
func GetDownloader(name string) (d Downloader) {
	downloadersLock.Lock()
	d, _ = downloaders[name]
	downloadersLock.Unlock()
	return
}

// Downloader is the interface to download the .torrent file.
type Downloader interface {
	// ResetURLRule allows that the downloader can reset the format of the url.
	// For the placeholder of the argument, please use '%s' instead.
	ResetURLRule(string)

	// SetTimeout resets the timeout the http request.
	SetTimeout(time.Duration)

	// SetProxy sets the http proxy.
	SetProxy(string) error

	// Download downloads the .torrent file about infohash.
	//
	// Notice: It should calculate the downloaded url by infohash based on
	// the url rule.
	Download(infohash string) ([]byte, error)
}
