package testing

import (
	"errors"
	"os"
	"path/filepath"
)

const LOCAL_TESTING_CONFIG = "../../config/test-local.json"
const DEFAULT_TESTING_CONFIG = "../../config/test.json"

func GetCurrentConfigPath() (string, error) {
	if path, found := os.LookupEnv("OAUTH_CONFIG_PATH"); found {
		return filepath.Abs(path)
	}

	if _, err := os.Stat(LOCAL_TESTING_CONFIG); err == nil {
		return filepath.Abs(LOCAL_TESTING_CONFIG)
	}

	if _, err := os.Stat(DEFAULT_TESTING_CONFIG); err == nil {
		return filepath.Abs(DEFAULT_TESTING_CONFIG)
	}

	return "", errors.New("config file setting is undefined.")
}
