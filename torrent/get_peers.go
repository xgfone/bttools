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
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/xgfone/bt/metainfo"
	"github.com/xgfone/bt/tracker"
)

func init() {
	RegisterCmd(&cli.Command{
		Name:      "getpeers",
		Usage:     "Get the peers of the torrent from the given tracker",
		ArgsUsage: "TORRENT_INFOHASH TRACER_URL",
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

			id := metainfo.NewRandomHash()
			infohash := metainfo.NewHashFromHexString(args[0])
			req := tracker.AnnounceRequest{InfoHash: infohash, PeerID: id}

			clients := make([]tracker.Client, 0)
			for _, t := range args[1:] {
				client, err := tracker.NewClient(t)
				if err != nil {
					return err
				}
				clients = append(clients, client)
			}

			timeout := ctx.Duration("timeout")
			c, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			wg := new(sync.WaitGroup)
			for _, client := range clients {
				wg.Add(1)
				go func(client tracker.Client) {
					defer wg.Done()
					resp, err := client.Announce(c, req)
					if err != nil {
						fmt.Printf("%s: %s\n", client.String(), err)
					} else {
						for _, addr := range resp.Addresses {
							fmt.Println(addr.String())
						}
					}
				}(client)
			}
			wg.Wait()

			return nil
		},
	})
}
