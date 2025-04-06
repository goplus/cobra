// Code generated by gop (Go+); DO NOT EDIT.

package main

import (
	"fmt"
	"github.com/goplus/cobra"
)

const _ = true

type mod struct {
	cobra.Command
	*App
}
type mod_init struct {
	cobra.Command
	*App
	LLGo    bool `flag:"llgo, val: true, usage: use LLGo as underlying compiler"`
	Verbose bool `flag:"verbose, short: v, usage: print verbose information"`
}
type version struct {
	cobra.Command
	*App
	Verbose bool `flag:"verbose, short: v, usage: print verbose information"`
}
type App struct {
	cobra.App
}
//line demo/godemo/mod_cmd.gox:1
func (this *mod) Main(_gop_arg0 string) {
	this.Command.Main(_gop_arg0)
//line demo/godemo/mod_cmd.gox:1:1
	this.Short("module maintenance")
//line demo/godemo/mod_cmd.gox:3:1
	this.Long(`Go mod provides access to operations on modules.

Note that support for modules is built into all the go commands,
not just 'go mod'. For example, day-to-day adding, removing, upgrading,
and downgrading of dependencies should be done using 'go get'.
See 'go help modules' for an overview of module functionality.
`)
//line demo/godemo/mod_cmd.gox:11:1
	this.Run(func() {
//line demo/godemo/mod_cmd.gox:12:1
		this.Help()
	})
}
func (this *mod) Classfname() string {
	return "mod"
}
//line demo/godemo/mod_init_cmd.gox:6
func (this *mod_init) Main(_gop_arg0 string) {
//line demo/godemo/mod_cmd.gox:11:1
	this.Command.Main(_gop_arg0)
//line demo/godemo/mod_init_cmd.gox:6:1
	this.Short("initialize new module in current directory")
//line demo/godemo/mod_init_cmd.gox:8:1
	this.Long(`Init initializes and writes a new go.mod file in the current directory, in
effect creating a new module rooted at the current directory. The go.mod file
must not already exist.

Init accepts one optional argument, the module path for the new module. If the
module path argument is omitted, init will attempt to infer the module path
using import comments in .go files, vendoring tool configuration files (like
Gopkg.lock), and the current directory (if in GOPATH).

If a configuration file for a vendoring tool is present, init will attempt to
import module requirements from it.

See https://golang.org/ref/mod#go-mod-init for more about 'go mod init'.
`)
//line demo/godemo/mod_init_cmd.gox:23:1
	this.Run(func() {
//line demo/godemo/mod_init_cmd.gox:24:1
		this.Printf("call go mod init: llgo=%v, verbose=%v\n", this.LLGo, this.Verbose)
	})
}
func (this *mod_init) Classfname() string {
	return "mod_init"
}
//line demo/godemo/version_cmd.gox:5
func (this *version) Main(_gop_arg0 string) {
//line demo/godemo/mod_init_cmd.gox:23:1
	this.Command.Main(_gop_arg0)
//line demo/godemo/version_cmd.gox:5:1
	this.Short("print Go version")
//line demo/godemo/version_cmd.gox:7:1
	this.Long(`Version prints the build information for Go binary files.

Go version reports the Go version used to build each of the named files.

If no files are named on the command line, go version prints its own
version information.

If a directory is named, go version walks that directory, recursively,
looking for recognized Go binaries and reporting their versions.
By default, go version does not report unrecognized files found
during a directory scan. The -v flag causes it to report unrecognized files.

The -m flag causes go version to print each file's embedded
module version information, when available. In the output, the module
information consists of multiple lines following the version line, each
indented by a leading tab character.

See also: go doc runtime/debug.BuildInfo.
`)
//line demo/godemo/version_cmd.gox:27:1
	this.Run(func() {
//line demo/godemo/version_cmd.gox:28:1
		fmt.Println("go1.0", "verbose:", this.Verbose)
	})
}
func (this *version) Classfname() string {
	return "version"
}
func (this *App) Main() {
	cobra.Gopt_App_Main(this, new(mod), new(mod_init), new(version))
}
func main() {
	new(App).Main()
}
