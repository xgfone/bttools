package torrent

import (
	"bytes"
	"fmt"
	"os"

	"github.com/urfave/cli"
	"github.com/xgfone/bttools/metainfo"
	"github.com/xgfone/go-tools/file"
)

func registerTorrentDumpCmd() {
	cmd := cli.Command{
		Name:    "dump",
		Aliases: []string{"d"},
		Usage:   "Dump the information of a .torrent file.",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "file, f",
				Usage: "Output the result into the file.",
			},
		},
		Action: cli.ActionFunc(dumpCmd),
	}

	commands = append(commands, cmd)
}

func dumpCmd(c *cli.Context) (err error) {
	args := c.Args()
	results := make([][]byte, len(args))
	for i, arg := range args {
		info, err := dumpMetaInfo(arg)
		if err != nil {
			return err
		}
		results[i] = info
	}

	buf := bytes.NewBuffer(nil)
	for i, info := range results {
		fmt.Fprintf(buf, "%s\n%s\n\n", args[i], info)
	}

	if c.String("file") == "" {
		_, err = fmt.Print(buf.Bytes())
	} else {
		_, err = file.WriteBytes(file.Abs(c.String("file")), buf.Bytes())
	}

	return
}

func dumpMetaInfo(file string) (info []byte, err error) {
	f, err := os.OpenFile(file, os.O_RDONLY, 0666)
	if err != nil {
		return
	}
	defer f.Close()

	m, err := metainfo.Parse(f)
	if err != nil {
		return
	}

	info, err = m.IndentedJSON("", "    ")
	return
}
