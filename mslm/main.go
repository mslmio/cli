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
	case cmd == "signup":
		err = cmdSignup()
	case cmd == "completion":
		err = cmdCompletion()
	case cmd == "version" || cmd == "vsn" || cmd == "v":
		err = cmdVersion()
	default:
		err = cmdDefault()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", progBase, err)
		os.Exit(1)
	}
}
