cobra - A Commander for modern Go+ CLI interactions
=====

[![Build Status](https://github.com/goplus/cobra/actions/workflows/go.yml/badge.svg)](https://github.com/goplus/cobra/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/goplus/cobra)](https://goreportcard.com/report/github.com/goplus/cobra)
[![GitHub release](https://img.shields.io/github/v/tag/goplus/cobra.svg?label=release)](https://github.com/goplus/cobra/releases)
[![GoDoc](https://pkg.go.dev/badge/github.com/goplus/cobra.svg)](https://pkg.go.dev/github.com/goplus/cobra)
[![Language](https://img.shields.io/badge/language-Go+-blue.svg)](https://github.com/goplus/gop)
<!--
[![Coverage Status](https://codecov.io/gh/goplus/cobra/branch/main/graph/badge.svg)](https://codecov.io/gh/goplus/cobra)
-->

The cobra classfile has the file suffix `_cmd.gox`.

## How to use in Go+

First let us initialize a `hellocli` project:

```sh
gop mod init hellocli
```

Then we have it reference the `cobra` classfile as the CLI framework:

```sh
gop get github.com/goplus/cobra@latest
```

Create a file named `version_cmd.gox` with the following content:

```go
run => {
	echo "command: version"
}
```

Execute the following commands:

```sh
gop mod tidy
gop install .
hellocli
```

You may get the following output:

```sh
Usage:
  hellocli [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version

Flags:
  -h, --help   help for hellocli

Use "hellocli [command] --help" for more information about a command.
```

## Command settings

Continue to modify the `version` command:

```go
// short sets the short description shown in the 'help' output.
short "print Go version"

// long sets the long message shown in the 'help <this-command>' output.
long `Version prints the build information for Go binary files.
`

run => {
	echo "command: version"
}
```

## Command flags

The `cobra` classfile uses tags of class fields to specify command flags.

```go
var (
	Verbose bool `flag:"verbose, short: v, usage: print verbose information"`
)

run => {
	echo "command: version", "verbose:", Verbose
}
```

## Subcommand

Create a file named `mod_cmd.gox` with the following content:

```go
run => {
	help
}
```

And create a file named `mod_init_cmd.gox` with the following content:

```go
run => {
	echo "subcommand: mod init"
}
```
