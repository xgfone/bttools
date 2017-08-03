package downloader

import "fmt"

func init() {
	RegisterDownloader("xunlei", NewXunLei())
}

const (
	xunleiURL = "http://bt.box.n0808.com/%s/%s/%s.torrent"
)

type xunlei struct {
	*HTTPClient
	url string
}

// NewXunLei returns a new Downloader based on XunLei network.
func NewXunLei() Downloader {
	return xunlei{HTTPClient: NewHTTPClient(), url: xunleiURL}
}

func (x xunlei) ResetURLRule(url string) {
	x.url = url
}

func (x xunlei) Download(infohash string) ([]byte, error) {
	url := fmt.Sprintf(x.url, infohash[:2], infohash[len(infohash)-2:], infohash)
	return x.Get(url)
}
