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
	"fmt"
	"path/filepath"
	"reflect"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/xgfone/bt/metainfo"
	"github.com/xgfone/bttools/pkg/helper"
)

func init() {
	registerCmd(&cli.Command{
		Name:      "showinfo",
		Usage:     "Print the metainfo information of the .torrent file and exit",
		ArgsUsage: "<TORRENT_FILES_PERTTERN> [TORRENT_FILES_PERTTERN ...]",
		Action: func(ctx *cli.Context) error {
			return printTorrentFiles(ctx.Args().Slice())
		},
	})
}

func printTorrentFiles(patterns []string) (err error) {
	// Get all the torrent files
	files := make([]string, 0, len(patterns))
	for _, s := range patterns {
		ss, err := filepath.Glob(s)
		if err != nil {
			return err
		}
		files = append(files, ss...)
	}

	for i, file := range files {
		if i > 0 {
			fmt.Println()
		}

		err = printTorrentFile(file)
		if err != nil {
			return
		}
	}

	return
}

func printTorrentFile(filename string) (err error) {
	mi, err := metainfo.LoadFromFile(filename)
	if err != nil {
		return fmt.Errorf("fail to load the torrent file '%s': %w", filename, err)
	}

	info, err := mi.Info()
	if err != nil {
		return fmt.Errorf("fail to decode metainfo of '%s': %s", filename, err)
	}

	infohash := mi.InfoHash()
	fmt.Printf("Magnet: %s\n", mi.Magnet(info.Name, infohash).String())
	fmt.Printf("InfoHash: %s\n", infohash)
	printValue("Encoding: ", mi.Encoding)
	printValue("CreatedBy: ", mi.CreatedBy)
	printValue("CreationDate: ", time.Unix(mi.CreationDate, 0).Format(time.RFC3339))
	printValue("Comment: ", mi.Comment)
	printTrackers(mi)
	printDHTNodes(mi)
	printURLList(mi)

	// Print the info part
	fmt.Println("Info:")
	fmt.Printf("    Name: %s\n", info.Name)
	fmt.Printf("    TotalLength: %s\n", helper.FormatSize(info.TotalLength()))
	fmt.Printf("    PieceLength: %s\n", helper.FormatSize(info.PieceLength))
	fmt.Printf("    PieceNumber: %d\n", info.CountPieces())
	printFiles("    ", info)

	return
}

func printValue(name string, v interface{}) {
	if !reflect.ValueOf(v).IsZero() {
		fmt.Printf("%s%v\n", name, v)
	}
}

func printTrackers(mi metainfo.MetaInfo) {
	announces := mi.Announces().Unique()
	if len(announces) == 0 {
		return
	}

	fmt.Println("Trackers:")
	for _, s := range announces {
		printValue("    ", s)
	}
}

func printDHTNodes(mi metainfo.MetaInfo) {
	if len(mi.Nodes) == 0 {
		return
	}

	fmt.Println("DHT Nodes:")
	for _, n := range mi.Nodes {
		printValue("    ", n)
	}
}

func printURLList(mi metainfo.MetaInfo) {
	if len(mi.URLList) == 0 {
		return
	}

	fmt.Println("WebSeed URLs:")
	for _, s := range mi.URLList {
		printValue("    ", s)
	}
}

func printFiles(prefix string, info metainfo.Info) {
	if !info.IsDir() {
		return
	}

	fmt.Printf("%sFiles:\n", prefix)
	for i, file := range info.AllFiles() {
		if i > 0 {
			fmt.Println()
		}

		fmt.Printf("%s    Path: %s\n", prefix, file.Path(info))
		fmt.Printf("%s    Length: %s\n", prefix, helper.FormatSize(file.Length))
	}
}
