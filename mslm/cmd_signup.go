package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/browser"
	bolt "go.etcd.io/bbolt"
)

type signupCli struct {
	SignupURL string `json:"signupURL"`
}

type apiKeyCli struct {
	ApiKey string `json:"apiKey"`
}

func cmdSignup() error {
	res, err := http.Get("http://localhost:1786/_/api/u/v1/signup/cli")
	if res.StatusCode == http.StatusTooManyRequests {
		return fmt.Errorf("too many requests")
	}
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Parse response.
	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	body := &signupCli{}
	err = json.Unmarshal(rawBody, body)
	if err != nil {
		return err
	}
	browser.OpenURL(body.SignupURL)
	fmt.Println("If the link does not open, please go to this link to get your API key:")
	fmt.Println("")
	fmt.Printf("%v\n", body.SignupURL)
	fmt.Println("")
	fmt.Println("Press [Enter] when done if not automatically detected.")

	// Retrieving CLI token from signup URL.
	parsedURL, err := url.Parse(body.SignupURL)
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

		res, err := http.Get("http://localhost:1786/_/api/u/v1/signup/check?cli_token=" + cliToken)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode == http.StatusOK {

			rawBody, err := io.ReadAll(res.Body)
			if err != nil {
				return err
			}
			body := &apiKeyCli{}
			err = json.Unmarshal(rawBody, body)
			if err != nil {
				return err
			}

			path, err := DbFilePath()
			if err != nil {
				return err
			}

			// Open the database.
			db, err := bolt.Open(path, 0600, nil)
			if err != nil {
				return err
			}
			defer db.Close()

			config, err := GetConfig(db)
			if err != nil {
				gConfig.ApiKey = body.ApiKey

				err = SaveConfig(gConfig, db)
				if err != nil {
					return err
				}
			} else {
				config.ApiKey = body.ApiKey
				err = SaveConfig(config, db)
				if err != nil {
					return err
				}
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
