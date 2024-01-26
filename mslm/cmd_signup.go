package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/mslmio/cli/lib/complete"
	"github.com/mslmio/cli/lib/complete/predict"
	"github.com/pkg/browser"
	"github.com/spf13/pflag"
)

// completionsSignup defines the completions for the "signup" command.
var completionsSignup = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
}

// printHelpSignup prints the help message for the "signup" command.
func printHelpSignup() {
	fmt.Printf(
		`Usage: %s signup [<opts>]

Description:
  Registers a new account on Mslm.

  The API key is automatically fetched after the signup flow is completed
  and when the email is verified.

Examples:
  # Signup.
  $ %[1]s signup

Options:
  General:
    --help, -h
      show help.
`, progBase)
}

type responseSignupUrl struct {
	Data signupCli `json:"data"`
}

type signupCli struct {
	SignupUrl string `json:"signup_url"`
}

type responseApiKey struct {
	Data apiKeyCli `json:"data"`
}

type apiKeyCli struct {
	ApiKey string `json:"api_key"`
}

func cmdSignup() error {
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.Parse()

	if fHelp {
		printHelpSignup()
		return nil
	}

	res, err := http.Get("https://mslm.io/_/api/u/v1/signup/cli")
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusTooManyRequests {
		return fmt.Errorf("too many requests")
	}
	defer res.Body.Close()

	// Parse response.
	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	body := &responseSignupUrl{}
	err = json.Unmarshal(rawBody, body)
	if err != nil {
		return err
	}

	browser.OpenURL(body.Data.SignupUrl)
	fmt.Println("If the link does not open, please go to this link to get your API key:")
	fmt.Println("")
	fmt.Printf("%v\n", body.Data.SignupUrl)
	fmt.Println("")
	fmt.Println("Press [Enter] when done if not automatically detected.")

	// Retrieving CLI token from signup URL.
	parsedURL, err := url.Parse(body.Data.SignupUrl)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return err
	}

	cliToken := parsedURL.Query().Get("cli_token")
	if cliToken == "" {
		fmt.Println("CLI token not found in URL")
		return err
	}

	// Check if signup flow is completed.
	maxAttempts := 200
	count := 0
	for {
		count++
		res, err := http.Get("https://mslm.io/_/api/u/v1/signup/cli/check?cli_token=" + cliToken)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode == http.StatusOK {
			rawBody, err := io.ReadAll(res.Body)
			if err != nil {
				return err
			}

			body := &responseApiKey{}
			err = json.Unmarshal(rawBody, body)
			if err != nil {
				return err
			}

			if err = SaveConfig("ApiKey", body.Data.ApiKey); err != nil {
				return err
			}

			fmt.Println("API Key fetched successfully.")
			break
		}

		if count == maxAttempts {
			if _, err := fmt.Scanln(); err != nil {
				return fmt.Errorf("%v", err)
			}
		}

		time.Sleep(time.Second)
	}

	return nil
}
