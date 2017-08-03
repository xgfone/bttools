package downloader

import (
	"bytes"
	"net/http"
	"net/url"
	"time"
)

// HTTPClient is a convenient http client to download .torrent file.
type HTTPClient struct {
	transport *http.Transport
	client    *http.Client
}

// NewHTTPClient returns a new HTTPClient.
func NewHTTPClient() *HTTPClient {
	transport := &http.Transport{}

	return &HTTPClient{
		transport: transport,
		client:    &http.Client{Transport: transport},
	}
}

// SetTimeout sets the timeout of the request and response.
func (h *HTTPClient) SetTimeout(timeout time.Duration) {
	h.client.Timeout = timeout
}

// SetProxy sets the http proxy server.
func (h *HTTPClient) SetProxy(rawurl string) error {
	_url, err := url.Parse(rawurl)
	if err != nil {
		return err
	}
	h.transport.Proxy = http.ProxyURL(_url)
	return nil
}

// Get downloads a .torrent file about infohash from url.
func (h *HTTPClient) Get(url string) (r []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := h.client.Do(req)
	if err != nil {
		return
	}

	if resp != nil {
		defer func() {
			// io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()
		}()

		if resp.StatusCode == 200 {
			buf := bytes.NewBuffer(nil)
			if _, err = buf.ReadFrom(resp.Body); err != nil {
				return
			}
			r = buf.Bytes()
		} else if resp.StatusCode == 404 {
			err = ErrNotFound
		} else {
			err = ErrServerRefuse
		}
	}

	return
}
