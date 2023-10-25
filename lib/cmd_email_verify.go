package lib

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/mslmio/sdk-go/email_verify"
	"github.com/mslmio/sdk-go/mslm"
	"github.com/spf13/pflag"
	"os"
)

// CmdEmailVerifyFlags defines the flags for the "email-verify" command.
type CmdEmailVerifyFlags struct {
	Help  bool
	JSON  bool
	CSV   bool
	YAML  bool
	token string
}

// Init initializes the common flags available to CmdEmailVerify with sensible
// defaults.
func (f *CmdEmailVerifyFlags) Init() {
	pflag.BoolVarP(
		&f.Help,
		"help", "h", false,
		"show help.",
	)
	pflag.BoolVarP(
		&f.JSON,
		"json", "j", false,
		"output as JSON.",
	)
	pflag.BoolVarP(
		&f.CSV,
		"csv", "c", false,
		"output as CSV.",
	)
	pflag.BoolVarP(
		&f.YAML,
		"yaml", "y", false,
		"output as YAML.",
	)
	pflag.StringVarP(
		&f.token,
		"token", "t", "",
		"the token to use.",
	)

}

// marshalToCSV is a custom CSV marshalling function to handle the mx field
// This is necessary because the mx field is a map[string]interface{}
// and the `OutputCSV` function does not handle this well
func marshalToCSV(data *email_verify.SingleVerifyResp) error {
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
		return err
	}

	// Set the JSON representation of mx in the records
	records[1][11] = string(mxJSON)

	// Write the CSV data
	err = csvWriter.WriteAll(records)
	if err != nil {
		return err
	}

	// Flush the writer to ensure data is written
	csvWriter.Flush()

	return nil
}

// CmdEmailVerify is the handler for the "email-verify" command.
func CmdEmailVerify(f CmdEmailVerifyFlags, args []string, printHelp func()) error {
	c := mslm.Init(f.token)

	if len(args) != 1 || f.Help {
		printHelp()
		return nil
	}

	resp, err := c.EmailVerify.SingleVerify(args[0])
	if err != nil {
		return fmt.Errorf("error verifying email: %w", err)
	}

	if f.YAML {
		err = OutputYAML(resp)
	} else if f.CSV {
		err = marshalToCSV(resp)
	} else {
		err = OutputJSON(resp)
	}

	return err
}
