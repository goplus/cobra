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

	// Two-level: "sandbox_list" → parent is sandbox, name is "list"
	parent, name := parentAndCmdName(root, cmds, "sandbox_list")
	if parent != sandbox.cobraCmd() || name != "list" {
		t.Fatalf("sandbox_list: got parent=%p name=%q, want sandbox cmd", parent, name)
	}

	// Three-level: "sandbox_template_list" → parent is sandbox_template, name is "list"
	parent, name = parentAndCmdName(root, cmds, "sandbox_template_list")
	if parent != sandboxTpl.cobraCmd() || name != "list" {
		t.Fatalf("sandbox_template_list: got parent=%p name=%q, want sandbox_template cmd", parent, name)
	}

	// Top-level: "version" → parent is root, name is "version"
	parent, name = parentAndCmdName(root, cmds, "version")
	if parent != &root.Command || name != "version" {
		t.Fatalf("version: got parent=%p name=%q, want root", parent, name)
	}

	// Intermediate fallback: "sandbox_templating_list" → rightmost split
	// "sandbox_templating" is not registered, but fallback matches "sandbox"
	// with the remainder "templating_list" as the command name.
	parent, name = parentAndCmdName(root, cmds, "sandbox_templating_list")
	if parent != sandbox.cobraCmd() || name != "templating_list" {
		t.Fatalf("sandbox_templating_list: got parent=%p name=%q, want sandbox cmd with name=templating_list", parent, name)
	}

	// Fallback: "unknown_sub_deep" with no matching parent → root, full name.
	// Also verifies that a warning is emitted for underscore-containing names
	// that find no parent, surfacing likely misregistrations.
	var logBuf bytes.Buffer
	origOutput := log.Writer()
	origFlags := log.Flags()
	defer func() {
		log.SetOutput(origOutput)
		log.SetFlags(origFlags)
	}()
	log.SetOutput(&logBuf)
	log.SetFlags(0)
	parent, name = parentAndCmdName(root, cmds, "unknown_sub_deep")
	if parent != &root.Command || name != "unknown_sub_deep" {
		t.Fatalf("unknown_sub_deep: got parent=%p name=%q, want root with full name", parent, name)
	}
	if !strings.Contains(logBuf.String(), `"unknown_sub_deep"`) {
		t.Fatalf("unknown_sub_deep: expected warning log, got %q", logBuf.String())
	}

	// Top-level with no underscore must not emit a warning.
	logBuf.Reset()
	parentAndCmdName(root, cmds, "version")
	if logBuf.Len() != 0 {
		t.Fatalf("version: expected no warning log, got %q", logBuf.String())
	}
}
