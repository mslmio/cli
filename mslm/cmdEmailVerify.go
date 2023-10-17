package main

import (
	"fmt"
	"github.com/mslmio/cli/lib"
	"github.com/mslmio/cli/lib/complete"
	"github.com/mslmio/cli/lib/complete/predict"
	"github.com/mslmio/sdk-go/mslm"
	"github.com/spf13/pflag"
)

var completionsEmailVerify = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-t":      predict.Nothing,
		"--token": predict.Nothing,
		"-h":      predict.Nothing,
		"--help":  predict.Nothing,
		"-y":      predict.Nothing,
		"--yaml":  predict.Nothing,
		"-j":      predict.Nothing,
		"--json":  predict.Nothing,
		"-c":      predict.Nothing,
		"--csv":   predict.Nothing,
	},
}

func printHelpEmailVerify() {
	fmt.Printf(
		`Usage: %s emailverify [<opts>] <email>

Examples:
  # Verify an email address.
  $ %[1]s emailverify <email>

Options:
  --token, -t <token>
	use <token> as API token.
  --yaml, -y
	output as YAML.
  --json, -j
	output as JSON.
  --csv, -c
	output as CSV.
  --help, -h
    show help.
`, progBase)
}

func cmdEmailVerify() error {
	var fYAML bool
	var fJSON bool
	var fCSV bool
	var token string

	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.BoolVarP(&fYAML, "yaml", "y", false, "output as YAML.")
	pflag.BoolVarP(&fJSON, "json", "j", false, "output as JSON.")
	pflag.BoolVarP(&fCSV, "csv", "c", false, "output as CSV.")
	pflag.StringVarP(&token, "token", "t", "", "the token to use.")
	pflag.Parse()

	c := mslm.Init(token)

	args := pflag.Args()[1:]

	if len(args) != 1 || fHelp {
		printHelpEmailVerify()
		return nil
	}

	resp, err := c.EmailVerify.SingleVerify(args[0])
	if err != nil {
		return fmt.Errorf("error verifying email: %w", err)
	}

	if fYAML {
		err = lib.OutputYAML(resp)
	} else if fCSV {
		err = lib.OutputCSV(resp)
	} else {
		err = lib.OutputJSON(resp)
	}

	if err != nil {
		return err
	}
	return nil
}
