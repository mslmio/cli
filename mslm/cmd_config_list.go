package main

import (
	"fmt"

	"github.com/mslmio/cli/lib/complete"
	"github.com/mslmio/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsConfigList = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
}

func printHelpConfigList() {
	fmt.Printf(
		`Usage: %s config list [<opts>]

Description:
  List all available configurations.

Options:
  --help, -h
    show help.

Configurations:
  api_key
    The API key used when querying the API.
`, progBase)
}

func cmdConfigList() error {
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.Parse()

	if fHelp {
		printHelpConfigList()
		return nil
	}

	printHelpConfigList()

	return nil
}
