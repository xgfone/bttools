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
	"context"
	"fmt"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/xgfone/bt/metainfo"
	"github.com/xgfone/bt/tracker"
)

func init() { registerCmd(getpeersCmd) }

var getpeersCmd = &cli.Command{
	Name:      "getpeers",
	Usage:     "Get the peers of the torrent from the tracker",
	ArgsUsage: "<TORRENT_FILE>",
	Flags: []cli.Flag{
		&cli.StringSliceFlag{
			Name:  "tracker",
			Value: cli.NewStringSlice(defaultTrackers...),
			Usage: "The URL of the default tracker",
		},
		&cli.DurationFlag{
			Name:  "timeout",
			Value: time.Second * 10,
			Usage: "The timeout to get peers from the tracker",
		},
	},

	Action: func(ctx *cli.Context) (err error) {
		args := ctx.Args().Slice()
		if len(args) != 1 {
			cli.ShowCommandHelpAndExit(ctx, "getpeers", 1)
		}

		mi, err := metainfo.LoadFromFile(args[0])
		if err != nil {
			return
		}

		info, err := mi.Info()
		if err != nil {
			return
		}

		trackers := mi.Announces().Unique()
		if len(trackers) == 0 {
			trackers = ctx.StringSlice("tracker")
		}

		timeout := ctx.Duration("timeout")
		c, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		infohash := mi.InfoHash()
		totalLength := info.TotalLength()
		nodeid := metainfo.NewRandomHash()
		for i, t := range trackers {
			resp, err := tracker.GetPeers(c, t, nodeid, infohash, totalLength)
			if err != nil {
				fmt.Printf("fail to try tracker '%s': %v\n", t, err)
				continue
			} else if len(resp.Addresses) == 0 {
				fmt.Printf("WARNING: no peers from tracker '%s': incomplete=%d, complete=%d\n",
					t, resp.Leechers, resp.Seeders)
				continue
			}

			if i > 0 {
				fmt.Println()
			}

			fmt.Printf("Tracker %s:\n", t)
			fmt.Printf("   Interval: %ds\n", resp.Interval)
			fmt.Printf("   Complete: %d\n", resp.Seeders)
			fmt.Printf("   Incomplete: %d\n", resp.Leechers)
			fmt.Printf("   Peers:\n")
			for _, addr := range resp.Addresses {
				fmt.Printf("        %s\n", addr.String())
			}
		}

		return
	},
}
