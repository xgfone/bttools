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

// Package command is used to collect all the commands.
package command

import (
	"github.com/urfave/cli/v2"
	"github.com/xgfone/bttools/pkg/command/torrent"
)

// Commands is the set of all the sub-commands.
var Commands []*cli.Command

// RegisterCmd registers the command as the sub-command of the root.
func RegisterCmd(cmd *cli.Command) {
	Commands = append(Commands, cmd)
}

func init() {
	RegisterCmd(torrent.Command)
}
