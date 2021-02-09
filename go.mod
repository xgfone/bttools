module github.com/xgfone/bttools

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/samuel/go-zookeeper v0.0.0-20201211165307-7117e9ea2414 // indirect
	github.com/urfave/cli/v2 v2.3.0
	github.com/xgfone/bt v0.4.0
	github.com/xgfone/gconf/v5 v5.0.0
	github.com/xgfone/goapp v0.17.0
	github.com/xgfone/gover v0.3.0
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777
	golang.org/x/sys v0.0.0-20210124154548-22da62e12c0c // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

go 1.11

replace github.com/xgfone/bt v0.4.0 => ../bt
