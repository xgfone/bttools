package torrent

import (
	"github.com/urfave/cli"
)

var (
	commands []cli.Command
)

func init() {
	commands = make([]cli.Command, 0)

	registerDownloadCmd()
	registerTorrentDumpCmd()
}

// GetTorrentCmd returns the torrent sub-command.
func GetTorrentCmd() cli.Command {
	return cli.Command{
		Name:        "torrent",
		Aliases:     []string{"d"},
		Usage:       "Handle a metainfo file.",
		Subcommands: commands,
	}
}
