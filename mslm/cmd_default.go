package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"os"
)

func printHelpDefault() {
	fmt.Printf(
		`Usage: %s <cmd> [<opts>] [<args>]

Commands:
  emailverify  verify an email address.
  completion   install or output shell auto-completion script.
  version      show current version.

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

func cmdDefault() (err error) {
	var fVsn bool

	pflag.BoolVarP(&fVsn, "version", "v", false, "print binary release number.")
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.Parse()

	if fHelp {
		printHelpDefault()
		return nil
	}

	if fVsn {
		fmt.Println(version)
		return nil
	}

	args := pflag.Args()
	if len(args) != 0 && args[0] != "-" {
		fmt.Printf("err: \"%s\" is not a command.\n", os.Args[1])
		fmt.Println()
	}

	printHelpDefault()
	return nil
}