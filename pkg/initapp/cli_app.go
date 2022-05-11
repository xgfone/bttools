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

package initapp

import (
	"fmt"
	stdlog "log"

	"github.com/urfave/cli/v2"
	"github.com/xgfone/go-atexit"
	"github.com/xgfone/go-log"
	"github.com/xgfone/go-log/writer"
	"github.com/xgfone/gover"
)

func init() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Fprintln(c.App.Writer, c.App.Version)
	}
}

// NewApp creates and returns a new cli.App.
func NewApp() *cli.App {
	app := cli.NewApp()
	app.Usage = "A BitTorrent Tools"
	app.Version = gover.Version
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "log.file",
			Usage: "The log file path",
		},
		&cli.StringFlag{
			Name:  "log.level",
			Value: "info",
			Usage: "The log level, such as debug, info, etc",
		},
	}
	app.Action = func(ctx *cli.Context) error {
		if logfile := ctx.String("log.file"); logfile != "" {
			file := log.FileWriter(logfile, "100M", 100)
			log.SetWriter(writer.SafeWriter(file))
			atexit.RegisterWithPriority(0, func() { file.Close() })
		}
		log.SetLevel(log.ParseLevel(ctx.String("log.level")))
		stdlog.SetOutput(log.DefaultLogger.WithDepth(2))
		return nil
	}
	app.After = func(ctx *cli.Context) error {
		atexit.Execute()
		return nil
	}
	return app
}
