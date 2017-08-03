package metainfo

import (
	"encoding/json"
	"fmt"
	"io"

	bencode "github.com/jackpal/bencode-go"
)

// FileDict into which torrent metafile is parsed and stored into.
type FileDict struct {
	Length int64    `json:"length"`
	Path   []string `json:"path"`
	Md5sum string   `json:"md5sum,omitempty"`
}

// InfoDict define
type InfoDict struct {
	FileDuration []int64 `json:"file-duration,omitempty" bencode:"file-duration"`
	FileMedia    []int64 `json:"file-media,omitempty" bencode:"file-media"`

	// Single file
	Name   string `json:"name"`
	Length int64  `json:"length"`
	Md5sum string `json:"md5sum,omitempty"`

	// Multiple files
	Files       []FileDict `json:"files,omitempty"`
	PieceLength int64      `json:"piece length,omitempty" bencode:"piece length"`
	Pieces      string     `json:"-"`
	Private     int64      `json:"-"`
}

// MetaInfo define
type MetaInfo struct {
	Info         InfoDict   `json:"info"`
	InfoHash     string     `json:"info hash,omitempty" bencode:"info hash"`
	Announce     string     `json:"announce"`
	AnnounceList [][]string `json:"announce-list" bencode:"announce-list"`
	CreationDate int64      `json:"creation date" bencode:"creation date"`
	Comment      string     `json:"comment,omitempty"`
	CreatedBy    string     `json:"created by" bencode:"created by"`
	Encoding     string     `json:"encoding"`
}

// Parse reads .torrent file, un-bencode it and load them into MetaInfo struct.
func (m *MetaInfo) Parse(r io.Reader) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("bencode unmarshal panic: %v", e)
		}
	}()

	// Decode bencoded metainfo file.
	// 丢弃piece部分，节省90%以上流量
	e := bencode.Unmarshal(r, m)
	if e != nil && e.Error() != "ignore piece" {
		return e
	}

	return
}

// GetPiecesList splits pieces string into an array of 20 byte SHA1 hashes.
func (m *MetaInfo) GetPiecesList() []string {
	var piecesList []string
	piecesLen := len(m.Info.Pieces)
	for i, j := 0, 0; i < piecesLen; i, j = i+20, j+1 {
		piecesList = append(piecesList, m.Info.Pieces[i:i+19])
	}
	return piecesList
}

// JSON converts the metainfo to json.
func (m *MetaInfo) JSON() ([]byte, error) {
	return json.Marshal(m)
}

// IndentedJSON converts the metainfo to the indented json.
func (m *MetaInfo) IndentedJSON(prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(m, prefix, indent)
}
