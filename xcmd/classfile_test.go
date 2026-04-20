/*
 * Copyright (c) 2025 The XGo Authors (xgo.dev). All rights reserved.
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
	"bytes"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/goplus/cobra"
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
	name, short, val, usage = parseFlag("times, usage: print verbose information")
	if name != "times" || short != "" || val != "" || usage != "print verbose information" {
		t.Fatal("parseFlag `times, usage: print verbose information`:", name, short, val, usage)
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

// mockCmd implements iCommandProto for testing parentAndCmdName.
type mockCmd struct {
	cmd   cobra.Command
	fname string
}

func (m *mockCmd) cobraCmd() *cobra.Command { return &m.cmd }
func (m *mockCmd) Main(name string)         { m.cmd.Use = name }
func (m *mockCmd) Classfname() string       { return m.fname }

func TestParentAndCmdName(t *testing.T) {
	root := &Command{}
	root.Use("app")

	sandbox := &mockCmd{fname: "sandbox"}
	sandboxTpl := &mockCmd{fname: "sandbox_template"}
	cmds := []iCommandProto{sandbox, sandboxTpl}

	// rootParent is a sentinel resolved at runtime to &root.Command; table
	// entries refer to it by comparing pointer identity in the assertion.
	rootParent := &root.Command

	tests := []struct {
		name       string
		fname      string
		wantParent *cobra.Command
		wantName   string
		// wantWarn requires a warning log to be emitted; wantNoWarn asserts
		// silence. Both default to false, which skips the log check.
		wantWarn   bool
		wantNoWarn bool
	}{
		{
			name:       "two-level match",
			fname:      "sandbox_list",
			wantParent: sandbox.cobraCmd(),
			wantName:   "list",
			wantNoWarn: true,
		},
		{
			name:       "three-level match",
			fname:      "sandbox_template_list",
			wantParent: sandboxTpl.cobraCmd(),
			wantName:   "list",
			wantNoWarn: true,
		},
		{
			name:       "top-level no underscore",
			fname:      "version",
			wantParent: rootParent,
			wantName:   "version",
			wantNoWarn: true,
		},
		{
			// Rightmost split "sandbox_templating" is not registered; the
			// algorithm falls back to "sandbox" with the remainder as name.
			name:       "intermediate fallback",
			fname:      "sandbox_templating_list",
			wantParent: sandbox.cobraCmd(),
			wantName:   "templating_list",
			wantNoWarn: true,
		},
		{
			// Underscore-containing name with no matching parent falls back
			// to root with the full fname; a warning surfaces the likely
			// misregistration.
			name:       "unmatched underscore name warns",
			fname:      "unknown_sub_deep",
			wantParent: rootParent,
			wantName:   "unknown_sub_deep",
			wantWarn:   true,
		},
		{
			// Trailing underscore yields an empty name on the first split;
			// that candidate must be skipped so the valid parent can match.
			name:       "trailing underscore skips empty name",
			fname:      "sandbox_",
			wantParent: rootParent,
			wantName:   "sandbox_",
			wantWarn:   true,
		},
	}

	origOutput := log.Writer()
	origFlags := log.Flags()
	defer func() {
		log.SetOutput(origOutput)
		log.SetFlags(origFlags)
	}()
	log.SetFlags(0)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var logBuf bytes.Buffer
			log.SetOutput(&logBuf)

			gotParent, gotName := parentAndCmdName(root, cmds, tc.fname)

			if gotParent != tc.wantParent {
				t.Fatalf("parent: got %p, want %p", gotParent, tc.wantParent)
			}
			if gotName != tc.wantName {
				t.Fatalf("name: got %q, want %q", gotName, tc.wantName)
			}
			switch {
			case tc.wantWarn && !strings.Contains(logBuf.String(), strconv.Quote(tc.fname)):
				t.Fatalf("expected warning log mentioning %q, got %q", tc.fname, logBuf.String())
			case tc.wantNoWarn && logBuf.Len() != 0:
				t.Fatalf("expected no warning log, got %q", logBuf.String())
			}
		})
	}
}
