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

// Package helper provides some convenient helpful functions.
package helper

import (
	"fmt"
	"strconv"
)

const (
	kb = 1024
	mb = kb * 1024
	gb = mb * 1024
	tb = gb * 1024
	pb = tb * 1024
)

// FormatSize formats the size followed by the unit.
func FormatSize(size int64) string {
	var base int64
	var unit string
	switch {
	case size < kb:
		return strconv.FormatInt(size, 10)

	case size < mb:
		base = kb
		unit = "KB"

	case size < gb:
		base = mb
		unit = "MB"

	case size < tb:
		base = gb
		unit = "GB"

	case size < pb:
		base = tb
		unit = "TB"

	default:
		base = pb
		unit = "PB"
	}

	if size%base == 0 {
		return fmt.Sprintf("%.0f%s", float64(size)/float64(base), unit)
	}
	return fmt.Sprintf("%.2f%s", float64(size)/float64(base), unit)
}
