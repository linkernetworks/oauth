package app

import (
	"os"
	"testing"

	"bitbucket.org/linkernetworks/aurora/src/oauth/app/config"
)

func TestNewServiceProvider(t *testing.T) {
	if configPath, found := os.LookupEnv("OAUTH_CONFIG_PATH"); found {
		appConfig := config.Read(configPath)
		as := NewServiceProvider(appConfig)
		if as.OAuthConfig.ExpiryDuration != 3600 {
			t.Error()
		}

		if as.OsinServer == nil {
			t.Errorf("osin.Service should not nil")
		}
	}
}
