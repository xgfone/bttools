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

// Package torrent provides the sub-commands about torrent.
package torrent

import "github.com/urfave/cli/v2"

// Command is the sub-command torrent.
var Command = &cli.Command{
	Name:  "torrent",
	Usage: "The torrent tools",
}

// registerCmd registers the command as the sub-command of the root.
func registerCmd(cmd *cli.Command) {
	Command.Subcommands = append(Command.Subcommands, cmd)
}
