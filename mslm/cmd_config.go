package main

import (
	"fmt"
	"os"

	"github.com/mslmio/cli/lib/complete"
	"github.com/mslmio/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsConfig = &complete.Command{
	Sub: map[string]*complete.Command{
		"list":  completionsConfigList,
		"get":   completionsConfigGet,
		"set":   completionsConfigSet,
		"clear": completionsConfigClear,
	},
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
}

func printHelpConfig() {
	fmt.Printf(
		`Usage: %s config <cmd> [<opts>]

Commands:
  list      prints all the configs available.
  get       get the value of a specified config.
  set       set the value of a specified config.
  clear     reset config.

Options:
  --help, -h
    show help.
`, progBase)
}

func configHelp() error {
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.Parse()

	if fHelp {
		printHelpConfig()
		return nil
	}

	printHelpConfig()
	return nil
}

func cmdConfig() error {
	var err error
	cmd := ""
	if len(os.Args) > 2 {
		cmd = os.Args[2]
	}

	switch {
	case cmd == "list":
		err = cmdConfigList()
	case cmd == "get":
		err = cmdConfigGet()
	case cmd == "set":
		err = cmdConfigSet()
	case cmd == "clear":
		err = cmdConfigClear()
	default:
		err = configHelp()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %v\n", err)
	}

	return nil
}
