package main

import (
	"fmt"
	"strings"

	"github.com/mslmio/cli/lib/complete"
	"github.com/mslmio/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsConfigSet = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
}

func printHelpConfigSet() {
	fmt.Printf(
		`Usage: %s config set [<key>=<value>...]

Description:
  Change a configuration.

Examples:
  $ %[1]s config set api_key=<your-key>

Options:
  --help, -h
    show help.

Configurations:
  api_key=<your-key>
    Save the API key for use when querying the API.
`, progBase)
}

func cmdConfigSet() error {
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.Parse()

	args := pflag.Args()[2:]

	if fHelp || len(args) < 1 {
		printHelpConfigSet()
		return nil
	}

	for _, arg := range args {
		confStr := strings.Split(arg, "=")
		key := strings.ToLower(confStr[0])
		if len(confStr) != 2 {
			if key == "api_key" {
				if err := UpdateConfigFieldAndSave("ApiKey", ""); err != nil {
					return err
				}
			}
			return fmt.Errorf("err: invalid key argument %s", key)
		}
		conf, err := GetConfig()
		if err != nil && conf == nil {
			return err
		} else {
			switch key {
			case "api_key":
				if err := UpdateConfigFieldAndSave("ApiKey", confStr[1]); err != nil {
					return err
				}
			default:
				return fmt.Errorf("err: invalid key argument %s", key)
			}
		}
	}

	return nil
}
