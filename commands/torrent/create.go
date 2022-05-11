// Copyright 2020~2022 xgfone
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
	"os"
	"path/filepath"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/xgfone/bt/bencode"
	"github.com/xgfone/bt/metainfo"
)

func init() { registerCmd(createCmd) }

var createCmd = &cli.Command{
	Name:      "create",
	Usage:     "Generate a .torrent file from a directory",
	ArgsUsage: "[TORRENT_DIRECTORY]",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "no-date",
			Aliases: []string{"d"},
			Usage:   "Leave the date field unset",
		},
		&cli.StringSliceFlag{
			Name:    "webseed",
			Aliases: []string{"w"},
			Usage:   "List of possible webseed URLs to use",
		},
		&cli.StringSliceFlag{
			Name:    "announce",
			Aliases: []string{"a"},
			Value:   cli.NewStringSlice("udp://tracker.openbittorrent.com:80/announce"),
			Usage:   "List of announce URLs to use",
		},
		&cli.StringFlag{
			Name:    "comment",
			Aliases: []string{"c"},
			Value:   "",
			Usage:   "Add a comment to the torrent file",
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
			Name:    "length",
			Aliases: []string{"l"},
			Value:   256,
			Usage:   "Piece length to use in kilobytes, default is 256. mktorrent syntax(powers of 2, 15-32) are also supported",
		},
	},
	Action: func(ctx *cli.Context) (err error) {
		dirs := ctx.Args().Slice()
		if len(dirs) == 0 {
			dirs = []string{"."}
		}

		var config CreateTorrentConfig
		config.Name = ctx.String("name")
		config.Output = ctx.String("output")
		config.NoDate = ctx.Bool("no-date")
		config.Comment = ctx.String("comment")
		config.WebSeeds = ctx.StringSlice("webseed")
		config.Announces = ctx.StringSlice("announce")

		if length := ctx.Int64("length"); length < 64 {
			config.PieceLength = length ^ 2
		} else {
			config.PieceLength = length * 1024
		}

		for _, dir := range dirs {
			config.RootDir, err = filepath.Abs(dir)
			if err != nil {
				return err
			}

			if config.Output == "" {
				name := config.Name
				if name == "" {
					name = filepath.Base(config.RootDir)
				}
				config.Output = name + ".torrent"
			}

			if err := CreateTorrent(config); err != nil {
				return err
			}
		}

		return nil
	},
}

// CreateTorrentConfig is the configuration information to create a .torrent file.
type CreateTorrentConfig struct {
	PieceLength int64
	Name        string
	RootDir     string
	Output      string
	Comment     string
	Announces   []string
	WebSeeds    []string
	NoDate      bool
}

// CreateTorrent creates a .torrent file.
func CreateTorrent(config CreateTorrentConfig) error {
	info, err := metainfo.NewInfoFromFilePath(config.RootDir, config.PieceLength)
	if err != nil {
		return err
	}

	if config.Name != "" {
		info.Name = config.Name
	}

	var mi metainfo.MetaInfo
	mi.Comment = config.Comment
	mi.InfoBytes, err = bencode.EncodeBytes(info)
	if err != nil {
		return err
	}

	switch len(config.Announces) {
	case 0:
	case 1:
		mi.Announce = config.Announces[0]
	default:
		mi.AnnounceList = metainfo.AnnounceList{config.Announces}
	}

	for _, seed := range config.WebSeeds {
		mi.URLList = append(mi.URLList, seed)
	}

	if !config.NoDate {
		mi.CreationDate = time.Now().Unix()
	}

	var out io.WriteCloser = os.Stdout
	if config.Output != "" {
		out, err = os.OpenFile(config.Output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return err
		}
		defer out.Close()
	}

	return mi.Write(out)
}
