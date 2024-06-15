package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/mslmio/cli/lib/complete"
	"github.com/mslmio/cli/lib/complete/predict"
	"github.com/spf13/pflag"
	"golang.org/x/term"
)

var completionsLogin = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-k":     predict.Nothing,
		"--key":  predict.Nothing,
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
}

func printHelpLogin() {
	fmt.Printf(
		`Usage: %s login [<opts>]

Description:
  Authenticate with an Mslm account using an API key.

Examples:
  # Login using an API key.
  $ %[1]s login --key <api-key>

  # Let the CLI prompt you for API key.
  $ %[1]s login

Options:
  --key <api-key>, -k <api-key>
    API key to login with.
    (this is potentially unsafe; let the CLI prompt you instead).
  --help, -h
    show help.
`, progBase)
}

func cmdLogin() error {
	var fKey string
	var fHelp bool

	pflag.StringVarP(&fKey, "key", "k", "", "the API key to save.")
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.Parse()

	if fHelp {
		printHelpLogin()
		return nil
	}

	// get args without subcommand.
	args := pflag.Args()[1:]

	// only key arg allowed.
	if len(args) > 1 {
		printHelpLogin()
		return nil
	}

	// allow only flag or arg for key but not both.
	if fKey != "" && len(args) > 0 {
		return errors.New("ambiguous key input source")
	}

	// get key, from flag or command line.
	// if it exists, we'll exit early as it's an implicit login.
	key := fKey
	if len(args) > 0 {
		key = args[0]
	}
	if key == "" {
		newKey, err := enterKey(key)
		if err != nil {
			return fmt.Errorf(err.Error())
		}

		key = newKey
	}

	if err := UpdateConfigFieldAndSave("ApiKey", key); err != nil {
		return err
	}

	fmt.Println("done")
	return nil
}

func enterKey(key string) (string, error) {
	for key == "" {
		fmt.Printf("Enter API key: ")
		keyraw, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			return "", err
		}

		key = string(keyraw[:])

		// exit if we have a key now.
		if key != "" {
			break
		}

		fmt.Println("please enter a key")
	}

	return key, nil
}
