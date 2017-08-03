package main

import (
	"os"

	"github.com/urfave/cli"
	"github.com/xgfone/bttools/cmds"
)

const version = "0.1.0"

func main() {
	app := cli.NewApp()
	app.Usage = "BT tool collection"
	app.Version = version
	app.EnableBashCompletion = true
	app.Commands = cmds.Commands
	cli.HandleExitCoder(app.Run(os.Args))
}
