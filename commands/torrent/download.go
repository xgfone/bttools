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
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/xgfone/bt/downloader"
	"github.com/xgfone/bt/metainfo"
	pp "github.com/xgfone/bt/peerprotocol"
	"github.com/xgfone/bt/tracker"
	"golang.org/x/net/context"
)

func init() {
	registerCmd(&cli.Command{
		Name:      "download",
		Usage:     "Download the file from the remote peers",
		ArgsUsage: "<TORRENT_FILE>",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:  "tracker",
				Usage: "The URL of the tracker.",
			},
			&cli.StringSliceFlag{
				Name:  "peer",
				Usage: "The address of the peer.",
			},
			&cli.StringFlag{
				Name:  "savedir",
				Usage: "The directory to save the downloaded file.",
			},
			&cli.UintFlag{
				Name:  "lastindex",
				Usage: "The index of the piece that is downloaded last",
			},
			&cli.UintFlag{
				Name:  "lastoffset",
				Usage: "The offset of the piece that is downloaded last"},
		},
		Action: func(ctx *cli.Context) error {
			args := ctx.Args().Slice()
			if len(args) == 0 {
				cli.ShowCommandHelpAndExit(ctx, "download", 0)
			}

			mi, err := metainfo.LoadFromFile(args[0])
			if err != nil {
				return err
			}

			infohash := mi.InfoHash()
			info, err := mi.Info()
			if err != nil {
				return err
			}

			id := metainfo.NewRandomHash()
			peers := ctx.StringSlice("peer")
			if len(peers) == 0 {
				ts := ctx.StringSlice("tracker")
				if peers, err = getPeersFromTrackers(id, mi, ts); err != nil {
					return err
				} else if len(peers) == 0 {
					return fmt.Errorf("no peers")
				}
			}

			savedir := ctx.String("savedir")
			w := metainfo.NewWriter(savedir, info, 0)
			defer w.Close()

			lastindex := ctx.Uint("lastindex")
			lastoffset := ctx.Uint("lastoffset")
			starttime := time.Now()
			dm := newDownloadManager(w, info, lastindex, lastoffset)
			downloadFile(peers, id, infohash, dm)
			if dm.IsFinished() {
				cost := time.Now().Sub(starttime)
				fmt.Printf("Finish downloading, cost %s\n", cost)

				r := metainfo.NewReader(savedir, info)
				defer r.Close()

				buf := make([]byte, info.PieceLength)
				for index := range info.Pieces {
					piece := info.Piece(index)
					n, err := r.ReadAt(buf[:piece.Length()], piece.Offset())
					if err != nil {
						fmt.Printf("fail to read the downloaded file: %s\n", err)
						return nil
					}

					if sha1.Sum(buf[:n]) != piece.Hash() {
						fmt.Printf("inconsistent hash of the piece at %d\n", index)
						return nil
					}
				}
				fmt.Println("The SHA1 checksum is OK")
			}
			return nil
		},
	})
}

func getPeersFromTrackers(id metainfo.Hash, mi metainfo.MetaInfo,
	trackers []string) (peers []string, err error) {
	if len(trackers) == 0 {
		trackers = mi.Announces().Unique()
		if len(trackers) == 0 {
			return nil, fmt.Errorf("no trackers")
		}
	}

	c, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	resps := tracker.GetPeers(c, id, mi.InfoHash(), trackers)
	for _, r := range resps {
		if r.Error != nil {
			log.Printf("tracker '%s' error: %s", r.Tracker, r.Error)
			continue
		}

		for _, addr := range r.Resp.Addresses {
			addrs := addr.String()
			nonexist := true
			for _, peer := range peers {
				if peer == addrs {
					nonexist = false
					break
				}
			}
			if nonexist {
				peers = append(peers, addrs)
			}
		}
	}

	return
}

func downloadFile(peers []string, id, infohash metainfo.Hash, dm *downloadManager) {
	peerslen := len(peers)
	for peerslen > 0 && !dm.IsFinished() {
		peerslen--
		peer := peers[peerslen]
		peers = peers[:peerslen]
		downloadFileFromPeer(peer, id, infohash, dm)
	}
}

func downloadFileFromPeer(peer string, id, infohash metainfo.Hash, dm *downloadManager) {
	pc, err := pp.NewPeerConnByDial(peer, id, infohash, time.Second*3)
	if err != nil {
		log.Printf("fail to dial '%s'", peer)
		return
	}
	defer pc.Close()

	dm.doing = false
	pc.Timeout = time.Second * 10
	if err = pc.Handshake(); err != nil {
		log.Printf("fail to handshake with '%s': %s", peer, err)
		return
	}

	bdh := downloader.NewBlockDownloadHandler(dm.writer.Info(), dm.OnBlock, dm.RequestBlock)
	if err = bdh.OnHandShake(pc); err != nil {
		log.Printf("handshake error with '%s': %s", peer, err)
		return
	}

	var msg pp.Message
	for !dm.IsFinished() {
		switch msg, err = pc.ReadMsg(); err {
		case nil:
			switch err = pc.HandleMessage(msg, bdh); err {
			case nil, pp.ErrChoked:
			default:
				log.Printf("fail to handle the msg from '%s': %s", peer, err)
				return
			}
		case io.EOF:
			log.Printf("got EOF from '%s'", peer)
			return
		default:
			log.Printf("fail to read the msg from '%s': %s", peer, err)
			return
		}
	}
}

func newDownloadManager(w metainfo.Writer, info metainfo.Info,
	lastindex, lastoffset uint) *downloadManager {
	length := info.Piece(int(lastindex)).Length() - int64(lastoffset)
	return &downloadManager{
		writer:  w,
		pindex:  uint32(lastindex),
		poffset: uint32(lastoffset),
		plength: length,
	}
}

type downloadManager struct {
	writer  metainfo.Writer
	pindex  uint32
	poffset uint32
	plength int64
	doing   bool
}

func (dm *downloadManager) IsFinished() bool {
	if dm.pindex >= uint32(dm.writer.Info().CountPieces()) {
		return true
	}
	return false
}

func (dm *downloadManager) OnBlock(index, offset uint32, b []byte) (err error) {
	if dm.pindex != index {
		return fmt.Errorf("inconsistent piece: old=%d, new=%d", dm.pindex, index)
	} else if dm.poffset != offset {
		return fmt.Errorf("inconsistent offset for piece '%d': old=%d, new=%d",
			index, dm.poffset, offset)
	}

	dm.doing = false
	n, err := dm.writer.WriteBlock(index, offset, b)
	if err == nil {
		dm.poffset = offset + uint32(n)
		dm.plength -= int64(n)
	}
	return
}

func (dm *downloadManager) RequestBlock(pc *pp.PeerConn) (err error) {
	if dm.doing {
		return
	}

	if dm.plength <= 0 {
		dm.pindex++
		if dm.IsFinished() {
			return
		}

		dm.poffset = 0
		dm.plength = dm.writer.Info().Piece(int(dm.pindex)).Length()
	}

	index := dm.pindex
	begin := dm.poffset
	length := uint32(downloader.BlockSize)
	if length > uint32(dm.plength) {
		length = uint32(dm.plength)
	}

	log.Printf("Request Block from '%s': index=%d, offset=%d, length=%d",
		pc.RemoteAddr().String(), index, begin, length)
	if err = pc.SendRequest(index, begin, length); err == nil {
		dm.doing = true
	}
	return
}
