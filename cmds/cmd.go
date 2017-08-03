package cmds

import (
	"github.com/urfave/cli"
	"github.com/xgfone/bttools/cmds/torrent"
)

var (
	// Commands is the commands of the program
	Commands []cli.Command
)

func init() {
	Commands = make([]cli.Command, 0)

	RegisterCommand(torrent.GetTorrentCmd())

}

// RegisterCommand registers a command.
func RegisterCommand(c cli.Command) {
	Commands = append(Commands, c)
}
