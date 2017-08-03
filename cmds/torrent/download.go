package torrent

import (
	"path/filepath"

	"github.com/urfave/cli"
	"github.com/xgfone/bttools/downloader"
	"github.com/xgfone/go-tools/file"
)

func registerDownloadCmd() {
	cmd := cli.Command{
		Name:    "download",
		Aliases: []string{"D"},
		Usage:   "Download the torrent about the infohash.",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "dir, d",
				Usage: "The direcotry to save the .torrent file.",
			},
			cli.StringFlag{
				Name:  "url, u",
				Usage: "Reset the url to download the .torrent file. If there is the placeholder of the argument, such as infohash, please use %s instead.",
			},
			cli.StringFlag{
				Name:  "source, s",
				Usage: "The torrent source from where to download the .torrent file. Support: xunlei, url, etc. If using url, you must give the option --url.",
				Value: "xunlei",
			},
			cli.StringFlag{
				Name:  "proxy, p",
				Usage: "Set the proxy of the http.",
			},
		},
		Action: cli.ActionFunc(downloadTorrentFile),
	}

	commands = append(commands, cmd)
}

func downloadTorrentFile(c *cli.Context) error {
	client := downloader.GetDownloader(c.String("source"))
	if client == nil {
		return downloader.ErrNotExist
	}

	if c.String("url") != "" {
		client.ResetURLRule(c.String("url"))
	}

	if c.String("proxy") != "" {
		if err := client.SetProxy(c.String("proxy")); err != nil {
			return err
		}
	}

	dir := file.Abs(c.String("dir"))
	for _, infohash := range c.Args() {
		content, err := client.Download(infohash)
		if err != nil {
			return err
		}
		_file := filepath.Join(dir, infohash+".torrent")
		file.WriteBytes(_file, content)
	}

	return nil
}
