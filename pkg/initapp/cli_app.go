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

	"github.com/urfave/cli/v2"
	"github.com/xgfone/gconf/v5"
	"github.com/xgfone/goapp/config"
	"github.com/xgfone/gover"
)

func init() {
	gconf.UnregisterOpts(gconf.ConfigFileOpt)
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Fprintln(c.App.Writer, c.App.Version)
	}
}

// NewApp creates and returns a new cli.App.
func NewApp() *cli.App {
	app := cli.NewApp()
	app.Usage = "A BitTorrent Tools"
	app.Version = gover.Version
	app.Before = func(*cli.Context) error { InitLogging(); return nil }
	app.Flags = config.ConvertOptsToCliFlags()
	return app
}
