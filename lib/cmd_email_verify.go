package lib

import (
	"fmt"
	"github.com/mslmio/sdk-go/mslm"
	"github.com/spf13/pflag"
)

// CmdEmailVerifyFlags defines the flags for the "email-verify" command.
type CmdEmailVerifyFlags struct {
	Help  bool
	Json  bool
	Csv   bool
	Yaml  bool
	token string
}

// Init initializes the common flags available to CmdEmailVerify.
func (f *CmdEmailVerifyFlags) Init() {
	pflag.BoolVarP(
		&f.Help,
		"help", "h", false,
		"show help.",
	)
	pflag.BoolVarP(
		&f.Json,
		"json", "j", false,
		"output as JSON.",
	)
	pflag.BoolVarP(
		&f.Csv,
		"csv", "c", false,
		"output as CSV.",
	)
	pflag.BoolVarP(
		&f.Yaml,
		"yaml", "y", false,
		"output as YAML.",
	)
	pflag.StringVarP(
		&f.token,
		"token", "t", "",
		"the token to use.",
	)

}

// CmdEmailVerify is the handler for the "emailverify" command.
func CmdEmailVerify(f CmdEmailVerifyFlags, args []string, printHelp func()) error {
	if len(args) != 1 || f.Help {
		printHelp()
		return nil
	}

	c := mslm.Init(f.token)

	resp, err := c.EmailVerify.SingleVerify(args[0])
	if err != nil {
		return fmt.Errorf("error verifying email: %w", err)
	}

	if f.Yaml {
		err = OutputYAML(resp)
	} else if f.Csv {
		err = OutputCSV(resp)
	} else {
		err = OutputJSON(resp)
	}

	return err
}
