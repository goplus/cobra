/*
 * Copyright (c) 2025 The GoPlus Authors (goplus.org). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package xcmd

import (
	"testing"
)

func TestParseFlag(t *testing.T) {
	name, short, val, usage := parseFlag("verbose, short: v, usage: print verbose information, etc.")
	if name != "verbose" || short != "v" || val != "" || usage != "print verbose information, etc." {
		t.Fatal("parseFlag `verbose, short: v, usage: print verbose information, etc.`:", name, short, val, usage)
	}
	name, short, val, usage = parseFlag("times, val: 3, usage:t,i,m,e,s")
	if name != "times" || short != "" || val != "3" || usage != "t,i,m,e,s" {
		t.Fatal("parseFlag `times, val: 3, usage:t,i,m,e,s`:", name, short, val, usage)
	}
	name, short, val, usage = parseFlag(",usage:t,i,m,e,s")
	if name != "" || short != "" || val != "" || usage != "t,i,m,e,s" {
		t.Fatal("parseFlag `,usage:t,i,m,e,s`:", name, short, val, usage)
	}
	name, short, val, usage = parseFlag("verbose")
	if name != "verbose" || short != "" || val != "" || usage != "" {
		t.Fatal("parseFlag `verbose`:", name, short, val, usage)
	}
	defer func() {
		if e := recover(); e == nil {
			t.Fatal("parseFlag: no error?")
		}
	}()
	parseFlag("verbose, unknown:")
}
