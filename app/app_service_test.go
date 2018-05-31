package app

import (
	"os"
	"testing"

	"github.com/linkernetworks/oauth/app/config"
)

func TestNewServiceProvider(t *testing.T) {
	// EXECUTOR_NUMBER is a jenkins environment variable
	if os.Getenv("EXECUTOR_NUMBER") != "" {
		t.Skip("Fix this for concurrent build")
	}

	if configPath, found := os.LookupEnv("OAUTH_CONFIG_PATH"); found {
		appConfig := config.Read(configPath)
		as := NewServiceProvider(*appConfig)
		if as.OAuthConfig.ExpiryDuration != 3600 {
			t.Error()
		}

		if as.OsinServer == nil {
			t.Errorf("osin.Service should not nil")
		}
	}
}
