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

package main

import (
	"github.com/urfave/cli/v2"
	"github.com/xgfone/gconf/v5"
	"github.com/xgfone/goapp"
)

var commands []*cli.Command

// RegisterCmd registers the command as the sub-command of the root.
func RegisterCmd(cmd *cli.Command) {
	commands = append(commands, cmd)
}

func initLogging() {
	logfile := gconf.GetString(goapp.LogOpts[0].Name)
	loglevel := gconf.GetString(goapp.LogOpts[1].Name)
	goapp.InitLogging(loglevel, logfile)
}

func main() {
	gconf.RegisterOpts(goapp.LogOpts...)
	app := cli.NewApp()
	app.Usage = "A BitTorrent Tools"
	app.Commands = commands
	app.RunAndExitOnError()
}
