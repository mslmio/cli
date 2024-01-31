package main

import (
	"fmt"

	"github.com/mslmio/cli/lib"
	"github.com/mslmio/cli/lib/complete"
	"github.com/mslmio/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

// completionsEmailVerify defines the completions for the "emailverify" command.
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

// printHelpEmailVerify prints the help message for the "emailverify" command.
func printHelpEmailVerify() {
	fmt.Printf(
		`Usage: %s emailverify [<opts>] <email>

Examples:
  # Verify an email address.
  %[1]s emailverify abc@example.com

Options:
  General:
    --help, -h
      show help.

  Formats:
    --token, -t <token>
      use <token> as API token.
    --yaml, -y
      output as YAML.
    --json, -j
      output as JSON.
    --csv, -c
      output as CSV.
`, progBase)
}

// cmdEmailVerify is the handler for the "emailverify" command.
func cmdEmailVerify() error {
	f := lib.CmdEmailVerifyFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdEmailVerify(f, pflag.Args()[1:], printHelpEmailVerify)
}
