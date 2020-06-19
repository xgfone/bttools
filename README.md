# bttools
A command tool collection about BitTorrent.

## 1 Install
```shell
$ go get -u github.com/xgfone/bttools
$ cd $GOPATH/src/github.com/xgfone/bttools
$ ./build.sh install
```

## 2 Commands

### 2.1 Command `torrent`

```shell
$ torrent -h
NAME:
   torrent - A BitTorrent Tools

USAGE:
   torrent [global options] command [command options] [arguments...]

VERSION:
   v0.2.0

COMMANDS:
   download   Download the file from the remote peers
   getpeers   Get the peers of the torrent from the given tracker
   printinfo  Print the metainfo of the torrent file
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

#### 2.2 Sub-Command `printinfo`
```shell
$ torrent printinfo ~/Downloads/gimp-2.10.18-setup-2.exe.torrent
MagNet: magnet:?xt=urn:btih:2aa1fff0d7ca65b149194ec42957d49cb27836db&dn=gimp-2.10.18-setup-2.exe&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969&tr=udp%3A%2F%2Ftracker.openbittorrent.com%3A80&tr=udp%3A%2F%2Fopen.demonii.com%3A1337
InfoHash: 2aa1fff0d7ca65b149194ec42957d49cb27836db
CreatedBy: mktorrent 1.1
CreationDate: 1585091362
Comment: GIMP 2.10.18 Installer for Microsoft Windows - 32 and 64 Bit - Update 2 - Fixes https://gitlab.gnome.org/GNOME/gegl/-/issues/231, which corrupted images once the swap file hit 2 GiB
Trackers:
    udp://tracker.coppersurfer.tk:6969
    udp://tracker.leechers-paradise.org:6969
    udp://tracker.openbittorrent.com:80
    udp://open.demonii.com:1337
URLs:
    http://gimper.net/downloads/pub/gimp/v2.10/windows/
    http://gimp.afri.cc/pub/gimp/v2.10/windows/
    http://de-mirror.gimper.net/pub/gimp/v2.10/windows/
    http://mirrors-br.go-parts.com/gimp/v2.10/windows/
    http://gimp.parentingamerica.com/v2.10/windows/
    http://gimp.raffsoftware.com/v2.10/windows/
    http://gimp.skazkaforyou.com/v2.10/windows/
    http://ftp.iut-bm.univ-fcomte.fr/gimp/v2.10/windows/
    http://mirror.ibcp.fr/pub/gimp/v2.10/windows/
    http://servingzone.com/mirrors/gimp/v2.10/windows/
    http://artfiles.org/gimp.org/v2.10/windows/
    http://gimp.cybermirror.org/v2.10/windows/
    http://ftp.fernuni-hagen.de/ftp-dir/pub/mirrors/www.gimp.org/v2.10/windows/
    http://ftp.gwdg.de/pub/misc/grafik/gimp/v2.10/windows/
    http://mirrors.zerg.biz/gimp/v2.10/windows/
    http://ftp.cc.uoc.gr/mirrors/gimp/v2.10/windows/
    http://ftp.heanet.ie/mirrors/ftp.gimp.org/pub/gimp/v2.10/windows/
    http://www.ring.gr.jp/pub/graphics/gimp/v2.10/windows/
    http://ftp.snt.utwente.nl/pub/software/gimp/gimp/v2.10/windows/
    http://ftp.nluug.nl/graphics/gimp/v2.10/windows/
    http://piotrkosoft.net/pub/mirrors/ftp.gimp.org/v2.10/windows/
    http://mirrors.dominios.pt/gimpv2.10/windows/
    http://mirrors.fe.up.pt/mirrors/ftp.gimp.org/v2.10/windows/
    http://mirrors.serverhost.ro/gimp/v2.10/windows/
    http://sunsite.rediris.es/mirror/gimp/v2.10/windows/
    http://ftp.sunet.se/pub/gnu/gimp/v2.10/windows/
    http://www.mirrorservice.org/sites/ftp.gimp.org/pub/gimp/v2.10/windows/
    http://gimp.cp-dev.com/v2.10/windows/
    http://mirror.hessmo.com/gimp/v2.10/windows/
    http://gimp.mirrors.hoobly.com/gimp/v2.10/windows/
    http://mirror.umd.edu/gimp/gimp/v2.10/windows/
    http://mirrors.zerg.biz/gimp/v2.10/windows/
    http://mirrors.xmission.com/gimp/v2.10/windows/
    http://download.gimp.org/pub/gimp/v2.10/windows/
Info:
    Name: gimp-2.10.18-setup-2.exe
    Length: 179640456
    PieceLength: 262144
    PieceNumber: 686
```

#### 2.3 Sub-Command `getpeers`
```shell
$ torrent getpeers 2aa1fff0d7ca65b149194ec42957d49cb27836db udp://tracker.leechers-paradise.org:6969
115.75.157.253:3869
95.56.201.138:35988
92.37.217.46:24212
89.188.125.68:21504
88.241.84.122:32006
86.38.251.108:39390
84.236.31.129:23758
82.98.56.241:11138
79.190.238.162:51413
78.132.204.142:35348
75.100.5.7:35142
62.11.154.212:18093
46.147.118.46:25856
36.77.93.100:49160
24.236.88.55:25285
......
```

#### 2.4 Sub-Command `download`
```shell
$ torrent.exe download ~/Downloads/gimp-2.10.18-setup-2.exe.torrent
2020/06/19 16:09:13 Request Block from '96.255.83.163:52742': index=0, offset=0, length=16384
2020/06/19 16:09:13 Request Block from '96.255.83.163:52742': index=0, offset=16384, length=16384
2020/06/19 16:09:13 Request Block from '96.255.83.163:52742': index=0, offset=32768, length=16384
2020/06/19 16:09:13 Request Block from '96.255.83.163:52742': index=0, offset=49152, length=16384
2020/06/19 16:09:14 Request Block from '96.255.83.163:52742': index=0, offset=65536, length=16384
2020/06/19 16:09:14 Request Block from '96.255.83.163:52742': index=0, offset=81920, length=16384
2020/06/19 16:09:14 Request Block from '96.255.83.163:52742': index=0, offset=98304, length=16384
2020/06/19 16:09:14 Request Block from '96.255.83.163:52742': index=0, offset=114688, length=16384
2020/06/19 16:09:15 Request Block from '96.255.83.163:52742': index=0, offset=131072, length=16384
2020/06/19 16:09:15 Request Block from '96.255.83.163:52742': index=0, offset=147456, length=16384
2020/06/19 16:09:15 Request Block from '96.255.83.163:52742': index=0, offset=163840, length=16384
2020/06/19 16:09:15 Request Block from '96.255.83.163:52742': index=0, offset=180224, length=16384
2020/06/19 16:09:16 Request Block from '96.255.83.163:52742': index=0, offset=196608, length=16384
2020/06/19 16:09:16 Request Block from '96.255.83.163:52742': index=0, offset=212992, length=16384
2020/06/19 16:09:16 Request Block from '96.255.83.163:52742': index=0, offset=229376, length=16384
2020/06/19 16:09:17 Request Block from '96.255.83.163:52742': index=0, offset=245760, length=16384
2020/06/19 16:09:17 Request Block from '96.255.83.163:52742': index=1, offset=0, length=16384
2020/06/19 16:09:17 Request Block from '96.255.83.163:52742': index=1, offset=16384, length=16384
......
2020/06/19 17:05:43 Request Block from '96.255.83.163:52742': index=685, offset=0, length=16384
2020/06/19 17:05:43 Request Block from '96.255.83.163:52742': index=685, offset=16384, length=16384
2020/06/19 17:05:44 Request Block from '96.255.83.163:52742': index=685, offset=32768, length=16384
2020/06/19 17:05:44 Request Block from '96.255.83.163:52742': index=685, offset=49152, length=16384
2020/06/19 17:05:44 Request Block from '96.255.83.163:52742': index=685, offset=65536, length=6280
Finish downloading, cost 56m31.739075s
The SHA1 checksum is OK
```

**Notice:** *The downloader downloads the file serially, not concurrently.*
