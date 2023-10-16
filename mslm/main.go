package main

import (
	"fmt"
	"os"
	"path/filepath"
)

var progBase = filepath.Base(os.Args[0])
var version = "1.0.0"

// global flags.
var fHelp bool

func printHelp() {
	fmt.Printf(
		`Usage: %s <cmd> [<opts>] [<args>]

Commands:
  emailverify   verify an email address.
  completion    install or output shell auto-completion script.
  version       show current version.

Options:
  General:
    --help, -h
    show help.
    --version, -v
    print binary release number.

  Formats:
    --json, -j
    output JSON format.
    --csv, -c
    output CSV format.
    --yaml, -y
    output YAML format.
`, progBase)
}

func main() {
	var err error
	var cmd string

	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	handleCompletions()

	switch {
	case cmd == "emailverify":
		err = cmdEmailVerify()
	case cmd == "completion":
		err = cmdCompletion()
	case cmd == "version" || cmd == "vsn" || cmd == "v":
		err = cmdVersion()
	default:
		printHelp()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", progBase, err)
		os.Exit(1)
	}
}
