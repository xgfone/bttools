package metainfo

import (
	"fmt"
	"strings"

	"github.com/xgfone/go-tools/str"
)

func init() {
	registerBaseUrn("sha1")
	registerBaseUrn("ed2k")
	registerBaseUrn("aich")
	registerBaseUrn("btih")
	registerBaseUrn("md5")
	registerBaseUrn("kzhash")
}

// Urn stands for a URN.
type Urn interface {
	Type() string
	Hash() string
	String() string
}

// UrnParser is a type to parse the URN.
type UrnParser func(string) (Urn, error)

type baseUrn struct {
	_type string
	hash  string
}

func newBaseUrn(t, h string) baseUrn {
	return baseUrn{_type: t, hash: h}
}

func (b baseUrn) Type() string {
	return b._type
}

func (b baseUrn) Hash() string {
	return b.hash
}

func (b baseUrn) String() string {
	return fmt.Sprintf("urn:%s:%s", b._type, b.hash)
}

func registerBaseUrn(typ string) {
	typ = strings.ToLower(typ)
	parseUrn := func(urn string) (Urn, error) {
		ss := strings.Split(urn, ":")
		if len(ss) != 3 || strings.ToLower(ss[0]) != "urn" || strings.ToLower(ss[1]) != typ {
			return nil, ErrInvalidFormat
		}
		return newBaseUrn(typ, ss[2]), nil
	}

	RegisterUrnParser(typ, parseUrn)
}

var (
	urnParsers = make(map[string]UrnParser)
)

// RegisterUrnParser registers a URN parser.
func RegisterUrnParser(typ string, p UrnParser) {
	urnParsers[strings.ToLower(typ)] = p
}

// ParseUrn parses a url to Urn.
func ParseUrn(urn string) (Urn, error) {
	if len(urn) < 6 || urn[:4] != "urn:" {
		return nil, ErrInvalidFormat
	}

	items := str.SplitStringN(urn, ":", 2)
	if len(items) != 3 || items[0] == "" || items[1] == "" || items[2] == "" {
		return nil, ErrInvalidFormat
	}

	if p, ok := urnParsers[strings.ToLower(items[1])]; ok {
		return p(urn)
	}
	return nil, fmt.Errorf("Cannot parse the URN: %s", items[1])
}
