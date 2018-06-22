package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

// Read return *AppConfig from json config file
func Read(configPath string) *AppConfig {
	configContent, err := ioutil.ReadFile(configPath)
	if err != nil {
		logrus.Fatalf("read config file %s error: %v\n", configPath, err)
	}

	var appConfig AppConfig
	err = json.Unmarshal(configContent, &appConfig)
	if err != nil {
		logrus.Fatalf("parse json to appconfig error: %v\n", err)
	}

	return &appConfig
}

// ReadCurrent read AppConfig from config file specific env
func ReadCurrent() *AppConfig {
	configPath := os.Getenv("OAUTH_CONFIG_PATH")
	if configPath == "" {
		configPath = "../../../config/default.json"
	}

	return Read(configPath)
}
