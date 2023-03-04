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
	"crypto/sha1"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/xgfone/bt/downloader"
	"github.com/xgfone/bt/metainfo"
	pp "github.com/xgfone/bt/peerprotocol"
	"github.com/xgfone/bt/tracker"
)

func init() { registerCmd(downloadCmd) }

var downloadCmd = &cli.Command{
	Name:      "download",
	Usage:     "Download the file from the remote peers by the .torrent file",
	ArgsUsage: "<TORRENT_FILE>",
	Flags: []cli.Flag{
		&cli.StringSliceFlag{
			Name:  "tracker",
			Value: cli.NewStringSlice(defaultTrackers...),
			Usage: "The URL of the default tracker",
		},
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"w"},
			Value:   "",
			Usage:   "The output directory to save the torrent files",
		},
		&cli.StringSliceFlag{
			Name:  "peer",
			Usage: "The address of the peers",
		},
	},

	Action: func(ctx *cli.Context) (err error) {
		args := ctx.Args().Slice()
		if len(args) != 1 {
			cli.ShowSubcommandHelpAndExit(ctx, 1)
		}

		mi, err := metainfo.LoadFromFile(args[0])
		if err != nil {
			return
		}

		info, err := mi.Info()
		if err != nil {
			return
		}

		infohash := mi.InfoHash()
		nodeid := metainfo.NewRandomHash()

		peers := ctx.StringSlice("peer")
		if len(peers) == 0 {
			trackers := mi.Announces().Unique()
			if len(trackers) == 0 {
				trackers = []string{ctx.String("tracker")}
			}

			peers, err = getPeersFromTrackers(infohash, nodeid, info, trackers)
			if err != nil {
				return
			}

			if len(peers) == 0 {
				fmt.Println("no found peers")
				return
			}
		}

		c, cancel := context.WithCancel(context.Background())
		defer cancel()

		tasks := make(chan metainfo.Piece, len(info.Pieces))
		for i := range info.Pieces {
			tasks <- info.Piece(i)
		}

		wg := new(sync.WaitGroup)
		wg.Add(len(peers))
		results := make(chan downloadResult)
		for _, peer := range peers {
			go startDownloadWorker(c, wg, nodeid, infohash, peer, tasks, results)
		}
		go func() {
			wg.Wait()
			fmt.Println("no peer workers and exit.")
			cancel()
		}()

		outputdir := ctx.String("output")
		w := metainfo.NewWriter(outputdir, info, 0600)
		defer w.Close()

		var donePieces int
	LOOP:
		for {
			select {
			case <-c.Done():
				break LOOP

			case result := <-results:
				_, err := w.WriteBlock(uint32(result.Index), 0, result.Data)
				if err != nil {
					fmt.Printf("fail to write the data of piece #%d\n", result.Index)
					break LOOP
				}

				donePieces++
				total := info.CountPieces()
				percent := 100 * float64(donePieces) / float64(total)
				fmt.Printf("(%0.2f) successfully download piece %d/%d \n", percent, donePieces, total)

				if donePieces >= total {
					fmt.Printf("the torrent file is saved at %s\n", filepath.Join(outputdir, info.Name))
					break LOOP
				}
			}
		}

		return
	},
}

func getPeersFromTrackers(infohash, nodeid metainfo.Hash, info metainfo.Info, trackers []string) (peers []string, err error) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	totalLength := info.TotalLength()
	for _, t := range trackers {
		resp, err := tracker.GetPeers(c, t, nodeid, infohash, totalLength)
		if err != nil {
			fmt.Printf("fail to try tracker '%s': %v\n", t, err)
			continue
		} else if len(resp.Addresses) == 0 {
			fmt.Printf("WARNING: no peers from tracker '%s': incomplete=%d, complete=%d\n",
				t, resp.Leechers, resp.Seeders)
			continue
		}

		for _, addr := range resp.Addresses {
			peers = append(peers, addr.String())
		}
	}
	return
}

type downloadResult struct {
	Index int
	Data  []byte
}

func startDownloadWorker(ctx context.Context, wg *sync.WaitGroup,
	nodeid metainfo.Hash, infohash metainfo.Hash, addr string,
	tasks chan metainfo.Piece, results chan<- downloadResult) {
	defer wg.Done()

	conn, err := pp.NewPeerConnByDial(addr, nodeid, infohash, time.Second*3)
	if err != nil {
		fmt.Printf("fail to connect to the peer '%s': %v\n", addr, err)
		return
	}
	defer conn.Close()

	if err = conn.Handshake(); err != nil {
		fmt.Printf("fail to handshake with the peer '%s': %v\n", addr, err)
		return
	}

	msg, err := conn.ReadMsg()
	if err != nil {
		fmt.Printf("fail to read the bitfield message from peer '%s': %v\n", addr, err)
		return
	}

	if msg.Type != pp.MTypeBitField {
		fmt.Printf("the first msg from peer '%s' after handshake is not bitfield: msgtype=%s\n", addr, msg.Type.String())
		return
	}
	conn.BitField = msg.BitField

	// Notice the peer that we can recieve the data.
	conn.SetUnchoked()
	conn.SetInterested()

	for {
		select {
		case <-ctx.Done():
			return

		case piece := <-tasks:
			index := uint32(piece.Index())

			// Check whether the peer has the piece.
			if !conn.PeerHasPiece(index) {
				tasks <- piece
				continue
			}

			// Try to download the piece.
			data := make([]byte, piece.Length())
			err = tryDownloadPiece(conn, 0, piece, data)
			if err != nil {
				fmt.Printf("fail to download piece #%d: %v\n", index, err)
				return
			}

			// Chech the hash of piece.
			if hash := sha1.Sum(data); piece.Hash() != hash {
				fmt.Printf("piece #%d expects hash %s, but got %x\n", index, piece.Hash().String(), hash)
				tasks <- piece
				continue
			}

			conn.SendHave(index)
			results <- downloadResult{Index: piece.Index(), Data: data}
		}
	}
}

func tryDownloadPiece(conn *pp.PeerConn, pieceNum int, piece metainfo.Piece, data []byte) (err error) {
	conn.SetTimeout(time.Minute)
	defer conn.SetTimeout(0)

	fmt.Printf("start to try download piece #%d\n", piece.Index())

	var downloaded, offset, backoff int
	for length := int(piece.Length()); downloaded < length; {
		if !conn.PeerChoked {
			if backoff < 5 && offset < length {
				blockSize := metainfo.BlockSize
				if rest := length - offset; rest < blockSize {
					blockSize = rest
				}

				err = conn.SendRequest(uint32(piece.Index()), uint32(offset), uint32(blockSize))
				if err != nil {
					return
				}

				backoff++
				offset += blockSize
			}
		}

		msg, err := conn.ReadMsg()
		if err != nil {
			return err
		} else if msg.Keepalive {
			continue
		}

		handler := downloader.NewBlockDownloadHandler(pieceNum, nil, func(i, o uint32, b []byte) error {
			if index := piece.Index(); index != int(i) {
				return fmt.Errorf("expect the data of index %d, but got index %d", index, i)
			} else if int(o) >= len(data) {
				return fmt.Errorf("got an unknown piece block data")
			}

			copy(data[o:], b)
			downloaded += len(b)
			backoff--
			return nil
		})

		err = conn.HandleMessage(msg, bep3Handler{handler})
		if err != nil {
			return err
		}
	}

	return
}

type bep3Handler struct {
	downloader.BlockDownloadHandler
}

func (h bep3Handler) OnMessage(conn *pp.PeerConn, msg pp.Message) error {
	fmt.Printf("recieved a unhandled torrent message: type=%s, raddr=%s", msg.Type, conn.RemoteAddr())
	return nil
}
