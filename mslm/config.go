package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"go.etcd.io/bbolt"
)

const (
	ConfigBucket = "config"
	ConfigKey    = "configKey"
)

var gConfig Config

type Config struct {
	ApiKey string `json:"apiKey"`
}

// gets the global config directory, creating it if necessary.
func getDbFileDir() (string, error) {
	cdir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	confDir := filepath.Join(cdir, "mslm")
	if err := os.MkdirAll(confDir, 0700); err != nil {
		return "", err
	}

	return confDir, nil
}

// returns the path to the config file.
func DbFilePath() (string, error) {
	confDir, err := getDbFileDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(confDir, "config.db"), nil
}

func SaveKeyInDB(apiKey string) error {
	path, err := DbFilePath()
	if err != nil {
		return err
	}

	// Open the database.
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	config, err := GetConfig(db)
	if err != nil {
		gConfig.ApiKey = apiKey

		err = SaveConfig(gConfig, db)
		if err != nil {
			return err
		}
	} else {
		config.ApiKey = apiKey
		err = SaveConfig(config, db)
		if err != nil {
			return err
		}
	}

	return nil
}

func SaveConfig(config Config, db *bbolt.DB) error {
	return db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(ConfigBucket))
		if err != nil {
			return err
		}

		// Marshal config struct to JSON
		configBytes, err := json.Marshal(config)
		if err != nil {
			return err
		}

		// Save serialized config to the bucket
		return bucket.Put([]byte(ConfigKey), configBytes)
	})
}

func GetConfig(db *bbolt.DB) (Config, error) {
	var config Config

	err := db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(ConfigBucket))
		if bucket == nil {
			return fmt.Errorf("%s bucket not found", ConfigBucket)
		}

		// Retrieve serialized config from the bucket
		configBytes := bucket.Get([]byte(ConfigKey))
		if configBytes == nil {
			return fmt.Errorf("%s key not found", ConfigKey)
		}

		// Unmarshal JSON to struct
		return json.Unmarshal(configBytes, &config)
	})

	if err != nil {
		return Config{}, err
	}

	return config, nil
}
