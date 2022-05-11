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
	"context"
	"fmt"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/xgfone/bt/metainfo"
	"github.com/xgfone/bt/tracker"
)

func init() {
	registerCmd(&cli.Command{
		Name:      "getpeers",
		Usage:     "Get the peers of the torrent from the given tracker",
		ArgsUsage: "<TORRENT_INFOHASH> <TRACER_URL> [TRACER_URL ...]",
		Flags: []cli.Flag{&cli.DurationFlag{
			Name:  "timeout",
			Value: time.Minute,
			Usage: "The timeout to get peers from the trackers.",
		}},
		Action: func(ctx *cli.Context) error {
			args := ctx.Args().Slice()
			if len(args) < 2 {
				cli.ShowCommandHelpAndExit(ctx, "getpeers", 0)
			}

			timeout := ctx.Duration("timeout")
			c, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			infohash := metainfo.NewHashFromHexString(args[0])
			resps := tracker.GetPeers(c, metainfo.Hash{}, infohash, args[1:])
			for _, r := range resps {
				if r.Error != nil {
					fmt.Printf("%s: %s\n", r.Tracker, r.Error.Error())
					continue
				}

				for _, addr := range r.Resp.Addresses {
					fmt.Println(addr.String())
				}
			}

			return nil
		},
	})
}
