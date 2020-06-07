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

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/urfave/cli/v2"
	"github.com/xgfone/bt/metainfo"
)

func init() { RegisterCmd(printCmd) }

var printCmd = &cli.Command{
	Name:  "print",
	Usage: "Print the metainfo of the torrent file",
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
	fmt.Printf("    TotalLength: %d\n", info.TotalLength())
	fmt.Printf("    PieceLength: %d\n", info.PieceLength)
	fmt.Printf("    PieceNumber: %d\n", info.CountPieces())
	fmt.Printf("    Files:\n")
	for _, file := range info.AllFiles() {
		fmt.Printf("        PathName: %s\n", file.Path(info))
		fmt.Printf("        Length: %d\n", file.Length)
		fmt.Println()
	}
}

func printValue(name string, v interface{}) {
	if !reflect.ValueOf(v).IsZero() {
		fmt.Printf("%s%v\n", name, v)
	}
}

func printTrackers(mi metainfo.MetaInfo) {
	if mi.Announce == "" && len(mi.AnnounceList) == 0 {
		return
	}
	fmt.Println("Trackers:")
	printValue("    ", mi.Announce)
	for _, ss := range mi.AnnounceList {
		for _, s := range ss {
			printValue("    ", s)
		}
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
