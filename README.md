# bttools
A tool collection about BitTorrent.

## Installation
```shell
$ go get -u github.com/xgfone/bttools
$ cd $GOPATH/src/github.com/xgfone/bttools && dep ensure
$ go install github.com/xgfone/bttools
```

## Usage

```shell
$ bttools.exe -h
NAME:
   bttools.exe - BT tool collection

USAGE:
   bttools.exe [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
     torrent, d  Handle a metainfo file.
     help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

### The sub-command `torrent` of `bttools`

```shell
$ bttools.exe torrent -h
NAME:
   bttools.exe torrent - Handle a metainfo file.

USAGE:
   bttools.exe torrent command [command options] [arguments...]

COMMANDS:
     download, D  Download the torrent about the infohash.
     dump, d      Dump the information of a .torrent file.

OPTIONS:
   --help, -h  show help
```

#### The sub-command `download` of `torrent`
```shell
$ bttools.exe torrent download -h
NAME:
   bttools.exe torrent download - Download the torrent about the infohash.

USAGE:
   bttools.exe torrent download [command options] [arguments...]

OPTIONS:
   --dir value, -d value     The direcotry to save the .torrent file.
   --url value, -u value     Reset the url to download the .torrent file.
                             If there is the placeholder of the argument,
                             such as infohash, please use %s instead.
   --source value, -s value  The torrent source from where to download the .torrent file.
                             Support: xunlei, url, etc. If using url, you must give the option
                             --url. (default: "xunlei")
   --proxy value, -p value   Set the proxy of the http.
```

#### The sub-command `dump` of `torrent`

```shell
$ bttools.exe torrent dump -h
NAME:
   bttools.exe torrent dump - Dump the information of a .torrent file.

USAGE:
   bttools.exe torrent dump [command options] [arguments...]

OPTIONS:
   --file value, -f value  Output the result into the file.
```


## Todo List

- [x] Download the .torrent file by the infohash.
- [x] Dump the information of a .torrent file to JSON.
- [ ] Download the files by the .torrent file.


## Platform

- [x] Windows
- [x] Mac OS X
- [x] Unix/Linux
