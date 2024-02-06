package main

import (
	"fmt"

	"github.com/mslmio/cli/lib/complete"
	"github.com/mslmio/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsLogout = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
}

func printHelpLogout() {
	fmt.Printf(
		`Usage: %s logout [<opts>]

Options:
  --help, -h
    show help.
`, progBase)
}

func cmdLogout() error {
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.Parse()

	if fHelp {
		printHelpLogout()
		return nil
	}

	// Check if not logged in.
	config, err := GetConfig()
	if err != nil && config == nil { // If db fails to open.
		return err
	} else if err != nil { // If db opens but no config exists.
		fmt.Println("not logged in")
		return nil
	} else { // If db opens and a config exists.
		if config.ApiKey == "" {
			fmt.Println("not logged in")
			return nil
		}

		if err = SaveFieldInConfig("ApiKey", ""); err != nil {
			return err
		}
	}

	fmt.Println("logged out")

	return nil
}
