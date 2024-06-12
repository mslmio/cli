package main

import (
	"fmt"
	"strings"

	"github.com/mslmio/cli/lib/complete"
	"github.com/mslmio/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsConfigGet = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
}

func printHelpConfigGet() {
	fmt.Printf(
		`Usage: %s config get [<key1>, <key2>...]

Description:
  Get the value of the specified configuation.

Examples:
  $ %[1]s config get api_key

Options:
  --help, -h
    show help.
`, progBase)
}

func cmdConfigGet() error {
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.Parse()

	args := pflag.Args()[2:]

	if fHelp || len(args) < 1 {
		printHelpConfigGet()
		return nil
	}

	validKeys := []string{"api_key"}

	processedKeys := make(map[string]bool)

	for _, arg := range args {
		key := strings.ToLower(arg)

		if processedKeys[key] {
			continue
		}

		// Check if the key is a valid key
		valid := false
		for _, validKey := range validKeys {
			if key == validKey {
				valid = true
				break
			}
		}

		if !valid {
			return fmt.Errorf("err: invalid key argument %s", key)
		}
		processedKeys[key] = true

		conf, err := GetConfig()
		if err != nil && conf == nil { // If db fails to open.
			return err
		} else if err != nil { // If no config exists.
			return err
		} else {
			switch key {
			case "api_key":
				if conf.ApiKey == "" {
					continue
				}
				fmt.Println(conf.ApiKey)
			}
		}

	}

	return nil
}
