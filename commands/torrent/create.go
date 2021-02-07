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
	"fmt"
	//  "bytes"
	//	"os"
	//	"path/filepath"
	//	"reflect"
	"log"

	"github.com/urfave/cli/v2"
	"github.com/xgfone/bt/metainfo"
)

func init() { registerCmd(createCmd) }

var createCmd = &cli.Command{
	Name:      "create",
	Usage:     "generate a .torrent from a directory.",
	ArgsUsage: "TBD",
	Flags: []cli.Flag{
		&cli.StringSliceFlag{
			Name:    "announce",
			Aliases: []string{"a"},
			Value:   cli.NewStringSlice("udp://tracker.openbittorrent.com:80/announce"),
			Usage:   "list of announce URLs to use",
		},
		&cli.StringSliceFlag{
			Name:    "webseed",
			Aliases: []string{"w"},
			Usage:   "list of possible web seed URLs to use",
		},
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Value:   "tools.torrent",
			Usage:   "path to output the .torrent file",
		},
		&cli.IntFlag{
			Name:    "length",
			Aliases: []string{"l"},
			Value:   256,
			Usage:   "path to output the .torrent file",
		},
	},
	Action: func(ctx *cli.Context) error {
		return CreateTorrent(ctx)
		return nil
	},
}

func CreateTorrent(ctx *cli.Context) error {
	dirs := ctx.Args().Slice()
	if len(dirs) > 1 {
		return fmt.Errorf("Input invalid, please use only one file or directory at ta time.")
	}
	info, err := metainfo.NewInfoFromFilePath(dirs[0], int64(ctx.Int("length")))
	if err != nil {
		return err
	}
	log.Println(info.CountPieces(), "pieces")
	var pieces []byte
	for n := 0; n < info.CountPieces(); n++ {
		pieces = append(pieces, info.Piece(n).Hash().Bytes()...)
	}
	meta := &metainfo.MetaInfo{
		InfoBytes: pieces,
	}
	infoc, err := meta.Info()
	if err != nil {
		return err
	}
	log.Println("Info generated", infoc.Name)
	return nil
}
