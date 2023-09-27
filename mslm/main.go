package main

import (
	"fmt"
	"os"
	"path/filepath"
)

var progBase = filepath.Base(os.Args[0])

// global flags.
var fHelp bool

func printHelp() {
	fmt.Printf(
		`Usage: %s <cmd> [<opts>] [<args>]
Examples:
	  # Verify an email address.
	  $ %[1]s emailverify <email>
Options:
	  --help, -h	
	    show help.
`, progBase)
}

// hello world
func main() {
	var err error
	var cmd string

	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	handleCompletions()

	switch {
	case cmd == "emailVerify":
		err = cmdEmailVerify()
	default:
		printHelp()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", progBase, err)
		os.Exit(1)
	}
}
