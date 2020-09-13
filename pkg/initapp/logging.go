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
	"github.com/xgfone/gconf/v5"
	"github.com/xgfone/goapp/log"
)

func init() {
	for i, opt := range log.LogOpts {
		opt.Aliases = nil
		log.LogOpts[i] = opt
	}
	gconf.RegisterOpts(log.LogOpts...)
}

// InitLogging initializes the logging.
func InitLogging() {
	logfile := gconf.GetString(log.LogOpts[0].Name)
	loglevel := gconf.GetString(log.LogOpts[1].Name)
	log.InitLogging(loglevel, logfile)
}
