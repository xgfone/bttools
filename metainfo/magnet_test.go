package metainfo

import (
	"fmt"
)

func ExampleMagnet_String() {
	murl := "magnet:?xt=urn:ed2k:354B15E68FB8F36D7CD88FF94116CDC1&xt=urn:btih:QHQXPYWMACKDWKP47RRVIV7VOURXFE5Q&xl=10826029&dn=mediawiki-1.15.1.tar.gz&tr=udp%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce&as=http%3A%2F%2Fdownload.wikimedia.org%2Fmediawiki%2F1.15%2Fmediawiki-1.15.1.tar.gz&xs=http%3A%2F%2Fcache.example.org%2FXRX2PEFXOOEJFRVUCX6HMZMKS5TWG4K5&xs=dchub://example.org"

	m := new(Magnet)
	if err := m.Parse(murl); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("First URN: type=%s, hash=%s\n", m.Xt[0].Type(), m.Xt[0].Hash())
	fmt.Printf("length=%d, filename=%s\n", m.Xl, m.Dn)
	fmt.Printf("url=%s\n", m.String())

	// Output:
	// First URN: type=ed2k, hash=354B15E68FB8F36D7CD88FF94116CDC1
	// length=10826029, filename=mediawiki-1.15.1.tar.gz
	// url=magnet:?xt=urn:ed2k:354B15E68FB8F36D7CD88FF94116CDC1&xt=urn:btih:QHQXPYWMACKDWKP47RRVIV7VOURXFE5Q&dn=mediawiki-1.15.1.tar.gz&xl=10826029&tr=udp%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce&as=http%3A%2F%2Fdownload.wikimedia.org%2Fmediawiki%2F1.15%2Fmediawiki-1.15.1.tar.gz&xs=http%3A%2F%2Fcache.example.org%2FXRX2PEFXOOEJFRVUCX6HMZMKS5TWG4K5&xs=dchub%3A%2F%2Fexample.org
}
