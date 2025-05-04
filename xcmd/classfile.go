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
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/goplus/cobra"
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

// FlagOff disables the flag parsing for this command.
func (p *Command) FlagOff() {
	p.Command.DisableFlagParsing = true
}

// Run sets the actual work function. Most commands will only implement this.
func (p *Command) Run__0(fn func()) {
	p.Command.Run = func(cmd *cobra.Command, args []string) {
		fn()
	}
}

// Run sets the actual work function. Most commands will only implement this.
func (p *Command) Run__1(fn func(args []string)) {
	p.Command.Run = func(cmd *cobra.Command, args []string) {
		fn(args)
	}
}

// -----------------------------------------------------------------------------

var (
	_ = (*Command).cobraCmd
)

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
}

type iCommandProto interface {
	cobraCmd() *cobra.Command
	Main(fname string)
	Classfname() string
}

// Gopt_App_Main is required by Go+ compiler as the entry of a Cobra project.
func Gopt_App_Main(app iAppProto, cmds ...iCommandProto) {
	projname := strings.TrimSuffix(filepath.Base(os.Args[0]), ".exe")
	root := app.initApp(projname)
	if me, ok := app.(interface{ MainEntry() }); ok {
		me.MainEntry()
	}
	for _, cmd := range cmds {
		v := reflect.ValueOf(cmd).Elem()
		self := cmd.cobraCmd()
		handleFlags(self, v)
		fname := cmd.Classfname()
		parent, name := parentAndCmdName(root, cmds, fname)
		cmd.Main(name)
		parent.AddCommand(self)
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

func handleFlags(self *cobra.Command, v reflect.Value) {
	t := v.Type()
	for i, n := 2, v.NumField(); i < n; i++ {
		tfld := t.Field(i)
		if flag := tfld.Tag.Get("flag"); flag != "" {
			flags := self.Flags()
			name, shorthand, val, usage := parseFlag(flag)
			fld := v.Field(i)
			p := fld.Addr().Interface()
			switch fld.Kind() {
			case reflect.Bool:
				flags.BoolVarP(p.(*bool), name, shorthand, parseBool(val, name), usage)
			case reflect.String:
				flags.StringVarP(p.(*string), name, shorthand, val, usage)
			case reflect.Int:
				flags.IntVarP(p.(*int), name, shorthand, parseInt(val, name), usage)
			case reflect.Float64:
				flags.Float64VarP(p.(*float64), name, shorthand, parseFloat(val, name), usage)
			default:
				log.Panicf("unsupported flag type `%s` for field `%s`", tfld.Type, tfld.Name)
			}
		}
	}
}

func parseFloat(val, name string) float64 {
	if val == "" {
		return 0
	}
	ret, e := strconv.ParseFloat(val, 64)
	if e != nil {
		log.Panicf("invalid value for flag `%s`: %v\n", name, val)
	}
	return ret
}

func parseInt(val, name string) int {
	if val == "" {
		return 0
	}
	ret, e := strconv.Atoi(val)
	if e != nil {
		log.Panicf("invalid value for flag `%s`: %v\n", name, val)
	}
	return ret
}

func parseBool(val, name string) bool {
	switch val {
	case "", "false":
		return false
	case "true":
		return true
	}
	log.Panicf("invalid value for flag `%s`: %v\n", name, val)
	return false
}

func parseFlag(flag string) (name, shorthand, val, usage string) {
	const spaces = " \t"
	var first = true
	var part, next string
	for flag != "" {
		pos := strings.IndexByte(flag, ',')
		if pos < 0 {
			part, next = flag, ""
		} else {
			part, next = flag[:pos], flag[pos+1:]
		}
		part := strings.TrimLeft(part, spaces)
		if first {
			name = part
			first = false
		} else if strings.HasPrefix(part, "short:") {
			shorthand = strings.TrimLeft(part[6:], spaces)
		} else if strings.HasPrefix(part, "val:") {
			val = strings.TrimLeft(part[4:], spaces)
		} else if strings.HasPrefix(part, "usage:") {
			usage = strings.TrimLeft(flag[len(flag)-len(part)-len(next)+5:], spaces)
			return
		} else {
			log.Panicf("invalid flag format `%s`", flag)
		}
		flag = next
	}
	return
}

// -----------------------------------------------------------------------------
