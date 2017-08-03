package downloader

import "fmt"

func init() {
	RegisterDownloader("url", NewXunLei())
}

type _url struct {
	*HTTPClient
	url string
}

// NewURL returns a new Downloader based on a url.
func NewURL() Downloader {
	return _url{HTTPClient: NewHTTPClient()}
}

func (u _url) ResetURLRule(url string) {
	u.url = url
}

func (u _url) Download(infohash string) ([]byte, error) {
	if u.url == "" {
		return nil, ErrEmptyURL
	}
	return u.Get(fmt.Sprintf(u.url, infohash))
}
