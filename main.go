// Copyright 2023 xgfone
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
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/xgfone/bttools/pkg/command"
	"github.com/xgfone/gover"
)

func printVersion(c *cli.Context) {
	fmt.Fprintln(c.App.Writer, c.App.Version)
}

func main() {
	app := cli.NewApp()
	app.Usage = "A BitTorrent Tools"
	app.Version = gover.Version
	app.Commands = command.Commands
	app.EnableBashCompletion = true
	cli.VersionPrinter = printVersion
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
