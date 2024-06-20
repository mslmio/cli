package main

import (
	"fmt"
	"strings"

	"github.com/mslmio/cli/lib/complete"
	"github.com/mslmio/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsConfigClear = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
		"-a":     predict.Nothing,
		"--all":  predict.Nothing,
	},
}

func printHelpConfigClear() {
	fmt.Printf(
		`Usage: %s config clear [<key1>, <key2>...]

Description:
  Reset a specified config or the entire config set.

Examples:
  # Clear a specified config
  $ %[1]s config clear api_key

  # Reset all configs
  $ %[1]s config clear --all

Options:
  --help, -h
    show help.
  --all, -a
    reset all configs.
`, progBase)
}

func cmdConfigClear() error {
	var fAll bool

	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.BoolVarP(&fAll, "all", "a", false, "reset all configs.")
	pflag.Parse()

	args := pflag.Args()[2:]

	if fHelp || (len(args) < 1 && !fAll) {
		printHelpConfigClear()
		return nil
	}

	if fAll {
		if err := ClearConfig(); err != nil {
			return err
		}

		fmt.Println("Cleared all configs.")

		return nil
	}

	conf, err := GetConfig()
	if err != nil && conf == nil { // If db fails to open.
		return err
	} else if err != nil { // If no config exists.
		return err
	}

	for _, arg := range args {
		key := strings.ToLower(arg)
		if key != "api_key" {
			return fmt.Errorf("err: invalid key argument %s", key)
		}

		switch key {
		case "api_key":
			conf.ApiKey = ""
		default:
			return fmt.Errorf("err: invalid key argument %s", key)
		}
	}

	if err := SaveConfig(*conf); err != nil {
		return err
	}

	fmt.Println("Cleared.")

	return nil
}
