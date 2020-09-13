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
	"os"
	"path/filepath"
	"reflect"

	"github.com/urfave/cli/v2"
	"github.com/xgfone/bt/metainfo"
)

func init() { registerCmd(printCmd) }

var printCmd = &cli.Command{
	Name:      "printinfo",
	Usage:     "Print the metainfo of the torrent file",
	ArgsUsage: "<TORRENT_FILES_PERTTERN> [TORRENT_FILES_PERTTERN ...]",
	Action: func(ctx *cli.Context) error {
		printTorrentFiles(ctx.Args().Slice())
		return nil
	},
}

func printTorrentFiles(patterns []string) {
	// Get all the torrent file
	files := make([]string, 0, len(patterns))
	for _, s := range patterns {
		ss, err := filepath.Glob(s)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		files = append(files, ss...)
	}

	// Print each torrent file
	for _, file := range files {
		printTorrentFile(file)
	}
}

func printTorrentFile(filename string) {
	mi, err := metainfo.LoadFromFile(filename)
	if err != nil {
		fmt.Printf("fail to load '%s': %s", filename, err)
		return
	}

	info, err := mi.Info()
	if err != nil {
		fmt.Printf("fail to decode info of '%s': %s", filename, err)
		return
	}

	infohash := mi.InfoHash()
	fmt.Printf("MagNet: %s\n", mi.Magnet(info.Name, infohash).String())
	fmt.Printf("InfoHash: %s\n", infohash)
	printValue("Encoding: ", mi.Encoding)
	printValue("CreatedBy: ", mi.CreatedBy)
	printValue("CreationDate: ", mi.CreationDate)
	printValue("Comment: ", mi.Comment)
	printTrackers(mi)
	printDHTNodes(mi)
	printURLList(mi)

	// Print the info part
	fmt.Println("Info:")
	fmt.Printf("    Name: %s\n", info.Name)
	printSingalFile("    ", info)
	fmt.Printf("    PieceLength: %d\n", info.PieceLength)
	fmt.Printf("    PieceNumber: %d\n", info.CountPieces())
	printMultiFiles("    ", info)
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

	fmt.Println("URLs:")
	for _, s := range mi.URLList {
		printValue("    ", s)
	}
}

func printSingalFile(prefix string, info metainfo.Info) {
	if !info.IsDir() {
		fmt.Printf("%sLength: %d\n", prefix, info.Length)
	}
}

func printMultiFiles(prefix string, info metainfo.Info) {
	if info.IsDir() {
		fmt.Printf("%sTotalLength: %d\n", prefix, info.TotalLength())
		fmt.Printf("%sFiles:\n", prefix)
		for _, file := range info.AllFiles() {
			fmt.Printf("%s    PathName: %s\n", prefix, file.Path(info))
			fmt.Printf("%s    Length: %d\n", prefix, file.Length)
			fmt.Println()
		}
	} else {
		fmt.Println()
	}
}
