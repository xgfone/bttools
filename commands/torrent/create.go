// Copyright 2020 xgfone
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package torrent

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/xgfone/bt/bencode"
	"github.com/xgfone/bt/metainfo"
)

func init() { registerCmd(createCmd) }

var createCmd = &cli.Command{
	Name:      "create",
	Usage:     "generate a .torrent from a directory.",
	ArgsUsage: "TBD",
	Flags: []cli.Flag{
		/*&cli.BoolFlag{
			Name:    "private",
			Aliases: []string{"p"},
			Usage:   "(placeholder) set the private flag (Currently unsupported)",
		},*/
		&cli.BoolFlag{
			Name:    "no-date",
			Aliases: []string{"d"},
			Usage:   "leave the date field unset",
		},
		&cli.StringSliceFlag{
			Name:    "webseed",
			Aliases: []string{"w"},
			Usage:   "list of possible web seed URLs to use",
		},
		&cli.StringSliceFlag{
			Name:    "announce",
			Aliases: []string{"a"},
			Value:   cli.NewStringSlice("udp://tracker.openbittorrent.com:80/announce"),
			Usage:   "list of announce URLs to use",
		},
		&cli.StringFlag{
			Name:    "comment",
			Aliases: []string{"c"},
			Value:   "",
			Usage:   "add a comment to the torrent file",
		},
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Value:   "",
			Usage:   "path to output the .torrent file, default is the file or directory name.torent",
		},
		&cli.StringFlag{
			Name:    "name",
			Aliases: []string{"n"},
			Value:   "",
			Usage:   "path to output the .torrent file, default is $name.torrent",
		},
		&cli.IntFlag{
			Name:    "length",
			Aliases: []string{"l"},
			Value:   256,
			Usage:   "piece length to use in kilobytes, default is 256. mktorrent syntax(powers of 2, 15-32) are also supported.",
		},
	},
	Action: func(ctx *cli.Context) error {
		piece_len, name, path, output, comment, announces, webseeds, nodate, private := torrentArgs(ctx)
		return CreateTorrent(piece_len, name, path, output, comment, announces, webseeds, nodate, private)
	},
}

func torrentArgs(ctx *cli.Context) (pl int64, n string, p string, o string, c string, an []string, ws []string, nd bool, pr bool) {
	dirs := ctx.Args().Slice()
	if len(dirs) > 1 {
		log.Println("Input invalid, please use only one file or directory at a time.")
		os.Exit(1)
	}
	p = dirs[0]

	n = ctx.String("name")
	if n == "" {
		n = p
	}

	o = ctx.String("output")
	if o == "" {
		o = n + ".torrent"
	}
	if ctx.Int("length") < 64 {
		pl = int64(ctx.Int("length") ^ 2)
	} else {
		pl = int64(ctx.Int("length") * 1000)
	}
	an = ctx.StringSlice("announce")
	ws = ctx.StringSlice("webseed")
	pr = false //ctx.Bool("private")
	nd = ctx.Bool("no-date")
	c = ctx.String("")

	return
}

func CreateTorrent(piece_len int64, name string, path string, output string, comment string, announces []string, webseeds []string, nodate bool, priv bool) error {

	info, err := metainfo.NewInfoFromFilePath(path, piece_len)
	if err != nil {
		return err
	}

	var mi metainfo.MetaInfo
	mi.InfoBytes, err = bencode.EncodeBytes(info)
	if err != nil {
		return err
	}

	// Set the announce information.
	switch len(announces) {
	case 0:
	case 1:
		mi.Announce = announces[0]
	default:
		mi.AnnounceList = metainfo.AnnounceList{announces}
	}

	switch len(webseeds) {
	case 0:
	default:
		mi.URLList = metainfo.URLList{}
		for _, seed := range webseeds {
			mi.URLList = append(mi.URLList, seed)
		}
	}

	if !nodate {
		mi.CreationDate = time.Now().Unix()
	}

	//	if private {
	//	  mi.Private = private
	//	}

	if comment != "" {
		mi.Comment = comment
	}

	var out io.WriteCloser = os.Stdout
	if o := output; o != "" {
		out, err = os.OpenFile(o, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return err
		}
		defer out.Close()
	}

	return mi.Write(out)
}
