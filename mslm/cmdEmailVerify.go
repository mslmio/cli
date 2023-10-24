package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/mslmio/cli/lib"
	"github.com/mslmio/cli/lib/complete"
	"github.com/mslmio/cli/lib/complete/predict"
	"github.com/mslmio/sdk-go/email_verify"
	"github.com/mslmio/sdk-go/mslm"
	"github.com/spf13/pflag"
	"os"
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

func marshalToCSV(data *email_verify.SingleVerifyResp) {
	// Create a CSV writer that writes to os.Stdout
	csvWriter := csv.NewWriter(os.Stdout)

	// Create a slice of string slices for the CSV data
	records := [][]string{
		{"email", "username", "domain", "malformed", "suggestion", "status", "has_mailbox", "accept_all", "disposable", "free", "role", "mx"},
		{data.Email, data.Username, data.Domain, fmt.Sprintf("%v", data.Malformed), data.Suggestion, data.Status, fmt.Sprintf("%v", data.HasMailbox), fmt.Sprintf("%v", data.AcceptAll), fmt.Sprintf("%v", data.Disposable), fmt.Sprintf("%v", data.Free), fmt.Sprintf("%v", data.Role), ""},
	}

	// Convert the mx field to a JSON string
	mxJSON, err := json.Marshal(data.Mx)
	if err != nil {
		fmt.Println("JSON Marshaling Error:", err)
		return
	}

	// Set the JSON representation of mx in the records
	records[1][11] = string(mxJSON)

	// Write the CSV data
	csvWriter.WriteAll(records)

	// Flush the writer to ensure data is written
	csvWriter.Flush()
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
		marshalToCSV(resp)
	} else {
		err = lib.OutputJSON(resp)
	}

	if err != nil {
		return err
	}
	return nil
}
