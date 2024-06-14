package main

import (
	"fmt"
	"reflect"

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
`, progBase)
}

func cmdConfigList() error {
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.Parse()

	if fHelp {
		printHelpConfigList()
		return nil
	}

	config, err := GetConfig()
	if err != nil {
		return fmt.Errorf("failed to get config: %v", err)
	}

	v := reflect.ValueOf(*config)
	t := v.Type()

	fmt.Println("Available configurations:")
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		fmt.Printf("  %s\n", jsonTag)
	}

	return nil
}
