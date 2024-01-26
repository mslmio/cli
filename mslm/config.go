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

type Config struct {
	ApiKey string `json:"apiKey"`
}

// Gets the global config directory, creating it if necessary.
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

// Returns the path to the config file.
func DbFilePath() (string, error) {
	confDir, err := getDbFileDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(confDir, "config.db"), nil
}

func UpdateConfigField(config *Config, fieldName string, newValue any) *Config {
	switch fieldName {
	case "ApiKey":
		config.ApiKey = newValue.(string)
	}
	return config
}

func SaveConfig(configName string, configValue any) error {
	// Check if a config already exists.
	conf, err := GetConfig()
	if err != nil && conf == nil { // If db fails to open.
		return err
	} else if err == nil && conf != nil { // If db opens and a config exists.
		conf = UpdateConfigField(conf, configName, configValue)
	} else { // If db opens but no config exists.
		conf = &Config{}
		conf = UpdateConfigField(conf, configName, configValue)
	}

	path, err := DbFilePath()
	if err != nil {
		return err
	}

	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(ConfigBucket))
		if err != nil {
			return err
		}

		// Marshal config struct to JSON.
		configBytes, err := json.Marshal(*conf)
		if err != nil {
			return err
		}

		// Save serialized config to the bucket.
		return bucket.Put([]byte(ConfigKey), configBytes)
	})
}

func GetConfig() (*Config, error) {
	var config Config

	path, err := DbFilePath()
	if err != nil {
		return nil, err
	}

	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(ConfigBucket))
		if bucket == nil {
			return fmt.Errorf("%s bucket not found", ConfigBucket)
		}

		// Retrieve serialized config from the bucket.
		configBytes := bucket.Get([]byte(ConfigKey))
		if configBytes == nil {
			return fmt.Errorf("%s key not found", ConfigKey)
		}

		// Unmarshal JSON to struct.
		return json.Unmarshal(configBytes, &config)
	})

	if err != nil {
		return &config, err
	}

	return &config, nil
}
