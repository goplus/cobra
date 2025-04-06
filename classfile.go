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

package cobra

import (
	"log"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
)

const (
	GopPackage = true
)

// -----------------------------------------------------------------------------

// Command is the worker class of Cobra classfile.
type Command struct {
	cobra.Command
}

func (p *Command) cobraCmd() *cobra.Command {
	return &p.Command
}

// Main is required by Go+ compiler as the entry of a Cobra command.
func (p *Command) Main(cmd string) {
	p.Command.Use = cmd
}

// Use sets the one-line usage message.
func (p *Command) Use(use string) {
	p.Command.Use = use
}

// Short sets the short description shown in the 'help' output.
func (p *Command) Short(s string) {
	p.Command.Short = s
}

// Long sets the long message shown in the 'help <this-command>' output.
func (p *Command) Long(l string) {
	p.Command.Long = l
}

// Run sets the actual work function. Most commands will only implement this.
func (p *Command) Run(fn func()) {
	p.Command.Run = func(cmd *cobra.Command, args []string) {
		fn()
	}
}

// -----------------------------------------------------------------------------

// App is the project class of Cobra classfile.
type App struct {
	Command
}

func (p *App) initApp(projname string) *Command {
	p.Use(projname)
	return &p.Command
}

type iAppProto interface {
	initApp(projname string) *Command
	Classprojname() string
}

type iCommandProto interface {
	cobraCmd() *cobra.Command
	Main(fname string)
	Classfname() string
}

// Gopt_App_Main is required by Go+ compiler as the entry of a Cobra project.
func Gopt_App_Main(app iAppProto, cmds ...iCommandProto) {
	root := app.initApp(app.Classprojname())
	if me, ok := app.(interface{ MainEntry() }); ok {
		me.MainEntry()
	}
	for _, cmd := range cmds {
		reflect.ValueOf(cmd).Elem().Field(1).Set(reflect.ValueOf(app)) // (*command).App = app
		fname := cmd.Classfname()
		parent, name := parentAndCmdName(root, cmds, fname)
		cmd.Main(name)
		parent.AddCommand(cmd.cobraCmd())
	}
	root.Execute()
}

func parentAndCmdName(root *Command, cmds []iCommandProto, fname string) (*cobra.Command, string) {
	pos := strings.IndexByte(fname, '_')
	if pos < 0 {
		return &root.Command, fname
	}
	subcmd, name := fname[:pos], fname[pos+1:]
	for _, v := range cmds {
		if v.Classfname() == subcmd {
			return v.cobraCmd(), name
		}
	}
	log.Panicf("Command `%s %s`: parent command not found", subcmd, name)
	return nil, ""
}

// -----------------------------------------------------------------------------
