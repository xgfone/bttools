# bttools

A command tool collection about BitTorrent. And you can consider it as the example of the development library [`bt`](https://github.com/xgfone/bt).

## 1 Install

```shell
$ make
```

## 2 Commands

### 2.1 Command `torrent`

```shell
$ bttools torrent -h
NAME:
   bttools torrent - The torrent tools

USAGE:
   bttools torrent command [command options] [arguments...]

COMMANDS:
   create    Generate a .torrent file from a file or directory
   download  Download the file from the remote peers by the .torrent file
   getpeers  Get the peers of the torrent from the tracker
   showinfo  Print the metainfo information of the .torrent file and exit
   help, h   Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```
