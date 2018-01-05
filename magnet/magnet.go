package magnet

import (
	"bytes"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/xgfone/go-tools/strings2"
)

var (
	// ErrInvalidFormat is returned when failed to parse the magnet or urn url.
	ErrInvalidFormat = fmt.Errorf("The format is invalid")
)

// Magnet is a struct of the magnet url.
//
// See https://en.wikipedia.org/wiki/Magnet_URI_scheme
type Magnet struct {
	// Display Name
	Dn string `json:"dn,omitempty"`

	// eXact Length, Size in bytes
	Xl int64 `json:"xl,omitempty"`

	// eXact Topic, URN
	Xt []Urn `json:"xt,omitempty"`

	// address TRacker, Tracker URL for BitTorrent downloads.
	Tr []string `json:"tr,omitemtpy"`

	// Keyword Topic, Key words for search
	Kt []string `json:"kt,omitempty"`

	// Acceptable Source, Web link to the file hash
	As []string `json:"as,omitempty"`

	// eXact Source, P2P link
	Xs []string `json:"xs,omitempty"`

	// Manifest Topic, Link to the metafile that contains a list of magnet.
	// It may be a http or URN link.
	Mt []string `json:"mt,omitempty"`
}

// Parse parses a magnet url.
//
// The url must start with "magnet:?".
func Parse(rawurl string) (m Magnet, err error) {
	if len(rawurl) < 12 || rawurl[:8] != "magnet:?" {
		err = ErrInvalidFormat
		return
	}

	var urn Urn
	items := strings.Split(rawurl[8:], "&")
	for _, item := range items {
		// two := strings.Split(item, "=")
		two := strings2.SplitStringN(item, "=", 1)
		if len(two) != 2 || two[1] == "" {
			continue
		}

		var value string
		if value, err = url.QueryUnescape(two[1]); err != nil {
			return
		}

		name := strings.ToLower(two[0])
		switch name {
		case "xt":
			if urn, err = ParseUrn(value); err != nil {
				return
			}
			if len(m.Xt) == 0 {
				m.Xt = []Urn{urn}
			} else {
				m.Xt = append(m.Xt, urn)
			}
		case "dn":
			m.Dn = value
		case "xl":
			if m.Xl, err = strconv.ParseInt(value, 10, 64); err != nil {
				return
			}
		case "tr":
			if len(m.Tr) == 0 {
				m.Tr = []string{value}
			} else {
				m.Tr = append(m.Tr, value)
			}
		case "kt":
			if len(m.Kt) == 0 {
				m.Kt = []string{value}
			} else {
				m.Kt = append(m.Kt, value)
			}
		case "as":
			if len(m.As) == 0 {
				m.As = []string{value}
			} else {
				m.As = append(m.As, value)
			}
		case "xs":
			if len(m.Xs) == 0 {
				m.Xs = []string{value}
			} else {
				m.Xs = append(m.Xs, value)
			}
		case "mt":
			if len(m.Mt) == 0 {
				m.Mt = []string{value}
			} else {
				m.Mt = append(m.Mt, value)
			}
		}
	}

	return
}

// String outputs the magnet url.
func (m Magnet) String() string {
	buf := bytes.NewBuffer(nil)

	for _, xt := range m.Xt {
		buf.WriteString(fmt.Sprintf("&xt=%s", xt.String()))
	}

	if m.Dn != "" {
		buf.WriteString("&dn=" + url.QueryEscape(m.Dn))
	}

	if m.Xl > 0 {
		buf.WriteString(fmt.Sprintf("&xl=%d", m.Xl))
	}

	m.outputStringList(buf, "tr", m.Tr)
	m.outputStringList(buf, "kt", m.Kt)
	m.outputStringList(buf, "as", m.As)
	m.outputStringList(buf, "xs", m.Xs)
	m.outputStringList(buf, "mt", m.Mt)

	_buf := buf.String()
	if len(_buf) == 0 {
		return ""
	}
	return "magnet:?" + _buf[1:]
}

func (m Magnet) outputStringList(buf *bytes.Buffer, name string, ss []string) {
	for _, s := range ss {
		buf.WriteString(fmt.Sprintf("&%s=%s", name, url.QueryEscape(s)))
	}
}
