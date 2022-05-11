// Copyright 2022 xgfone
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

package peer

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// Command is used to start a BT peer server.
var Command = &cli.Command{
	Name:  "peer",
	Usage: "A BT peer server",
	Flags: []cli.Flag{},
	Action: func(ctx *cli.Context) error {
		fmt.Println("Unimplemented")
		return nil
	},
}
