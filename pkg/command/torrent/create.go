// Copyright 2023 xgfone
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
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/xgfone/bt/bencode"
	"github.com/xgfone/bt/metainfo"
)

func init() { registerCmd(createCmd) }

var createCmd = &cli.Command{
	Name:      "create",
	Usage:     "Generate a .torrent file from a file or directory",
	ArgsUsage: "[<FILE | DIRECTORY> ...]",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "no-date",
			Usage: "Leave the date field unset",
		},
		&cli.StringSliceFlag{
			Name:  "webseed",
			Usage: "List of possible webseed URLs to use",
		},
		&cli.StringSliceFlag{
			Name:  "announce",
			Value: cli.NewStringSlice("udp://tracker.openbittorrent.com:80/announce"),
			Usage: "List of announce URLs to use",
		},
		&cli.StringSliceFlag{
			Name:  "dhtnode",
			Usage: "List of the DHT nodes",
		},
		&cli.StringFlag{
			Name:    "comment",
			Aliases: []string{"c"},
			Value:   "",
			Usage:   "Add a comment to the torrent file",
		},
		&cli.StringFlag{
			Name:  "creator",
			Value: "",
			Usage: "The creator to create the torrent",
		},
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Value:   "",
			Usage:   "The path of .torrent file to be output",
		},
		&cli.StringFlag{
			Name:    "name",
			Aliases: []string{"n"},
			Value:   "",
			Usage:   "The name of the torrent",
		},
		&cli.Int64Flag{
			Name:  "length",
			Value: 256,
			Usage: "Piece length to use in kilobytes, default is 256. mktorrent syntax(powers of 2, 15-32) are also supported",
		},
	},
	Action: func(ctx *cli.Context) (err error) {
		dirs := ctx.Args().Slice()
		if len(dirs) == 0 {
			dirs = []string{"."}
		}

		mi := metainfo.MetaInfo{
			Comment:   ctx.String("comment"),
			CreatedBy: ctx.String("creator"),
		}

		if !ctx.Bool("no-date") {
			mi.CreationDate = time.Now().Unix()
		}

		switch announces := ctx.StringSlice("announce"); len(announces) {
		case 0:
		case 1:
			mi.Announce = announces[0]
		default:
			mi.AnnounceList = metainfo.AnnounceList{announces}
		}

		for _, seed := range ctx.StringSlice("webseed") {
			mi.URLList = append(mi.URLList, seed)
		}
		for _, node := range ctx.StringSlice("dhtnode") {
			hostaddr, err := metainfo.ParseHostAddr(node)
			if err != nil {
				return err
			}
			mi.Nodes = append(mi.Nodes, hostaddr)
		}

		var pieceLength int64
		if length := ctx.Int64("length"); length < 64 {
			pieceLength = length ^ 2
		} else {
			pieceLength = length * 1024
		}

		name := ctx.String("name")
		if len(dirs) > 1 {
			name = ""
		}

		mis := make([]metainfo.MetaInfo, len(dirs))
		for i, dir := range dirs {
			mis[i] = mi
			err = updateMetaInfo(&mis[i], dir, name, pieceLength)
			if err != nil {
				return
			}
		}

		output := ctx.String("output")
		for _, mi := range mis {
			filename := filepath.Join(output, mi.InfoHash().String()+".torrent")
			err = outputMetaInfo(mi, filename)
			if err != nil {
				return
			}

			fmt.Printf("successfully create the .torrent file to %s\n", filename)
		}

		return nil
	},
}

func updateMetaInfo(mi *metainfo.MetaInfo, dir, name string, pieceLength int64) (err error) {
	dir, err = filepath.Abs(dir)
	if err != nil {
		return
	}

	info, err := metainfo.NewInfoFromFilePath(dir, pieceLength)
	if err != nil {
		return err
	}

	if name != "" {
		info.Name = name
	}

	mi.InfoBytes, err = bencode.EncodeBytes(info)
	return
}

func outputMetaInfo(mi metainfo.MetaInfo, filename string) (err error) {
	buf := bytes.NewBuffer(make([]byte, 0, 16*32))
	if err = mi.Write(buf); err == nil {
		err = ioutil.WriteFile(filename, buf.Bytes(), 0600)
	}
	return
}
